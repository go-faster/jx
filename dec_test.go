package jx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func decoderOnlyError[T any](fn func(*Decoder) (T, error)) func(*Decoder) error {
	return func(d *Decoder) error {
		_, err := fn(d)
		return err
	}
}

// testBufferReader runs the given test function with various decoding modes (buffered, streaming, etc.).
func testBufferReader(input string, cb func(t *testing.T, d *Decoder)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Buffer", func(t *testing.T) {
			cb(t, DecodeStr(input))
		})

		t.Run("Reader", func(t *testing.T) {
			r := strings.NewReader(input)
			cb(t, Decode(r, 512))
		})

		t.Run("OneByteReader", func(t *testing.T) {
			r := strings.NewReader(input)
			obr := iotest.OneByteReader(r)
			cb(t, Decode(obr, 512))
		})

		t.Run("DataErrReader", func(t *testing.T) {
			r := strings.NewReader(input)
			obr := iotest.DataErrReader(r)
			cb(t, Decode(obr, 512))
		})
	}
}

func createTestCase(input string, cb func(t *testing.T, d *Decoder) error) func(t *testing.T) {
	valid := json.Valid([]byte(input))
	return testBufferReader(input, func(t *testing.T, d *Decoder) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("Input: %q", input)
			}
		})

		err := cb(t, d)
		if valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	})
}

func runTestCases(t *testing.T, cases []string, cb func(t *testing.T, d *Decoder) error) {
	for i, input := range cases {
		t.Run(fmt.Sprintf("Test%d", i), createTestCase(input, cb))
	}
}

func TestType_String(t *testing.T) {
	met := map[string]bool{}
	for i := Invalid; i <= Object+1; i++ {
		s := i.String()
		if s == "" {
			t.Error("blank")
		}
		if met[s] {
			t.Errorf("met %s", s)
		}
		met[s] = true
	}
	if len(met) != 8 {
		t.Error("unexpected met types")
	}
}

func TestDecoder_Reset(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		var d Decoder
		d.ResetBytes([]byte{})
		d.Reset(bytes.NewBufferString(`true`))
		v, err := d.Bool()
		require.NoError(t, err)
		require.True(t, v)
	})
	t.Run("ZeroLen", func(t *testing.T) {
		var d Decoder
		d.ResetBytes(make([]byte, 0, 100))
		d.Reset(bytes.NewBufferString(`true`))
		v, err := d.Bool()
		require.NoError(t, err)
		require.True(t, v)
	})
	t.Run("ZeroReader", func(t *testing.T) {
		v, err := Decode(bytes.NewBufferString(`true`), 0).Bool()
		require.NoError(t, err)
		require.True(t, v)
	})
}

func TestDecoderNegativeDepth(t *testing.T) {
	require.ErrorIs(t, GetDecoder().decDepth(), errNegativeDepth)
}

func TestDecoderBig(t *testing.T) {
	t.Run("Float", func(t *testing.T) {
		n, err := DecodeStr(`10000005125004315341545.1234215253`).BigFloat()
		require.NoError(t, err)
		require.NotNil(t, n)
		t.Run("Negative", func(t *testing.T) {
			t.Run("Bad", func(t *testing.T) {
				for _, s := range []string{
					"",
				} {
					_, err := DecodeStr(s).BigFloat()
					require.Error(t, err)
				}
			})
			t.Run("Reader", func(t *testing.T) {
				_, err := Decode(errReader{}, 10).BigFloat()
				require.Error(t, err)
			})
		})
	})
	t.Run("Int", func(t *testing.T) {
		n, err := DecodeStr(`10000005125004315341545`).BigInt()
		require.NoError(t, err)
		require.NotNil(t, n)
		t.Run("Negative", func(t *testing.T) {
			t.Run("Bad", func(t *testing.T) {
				for _, s := range []string{
					"",
				} {
					_, err := DecodeStr(s).BigInt()
					require.Error(t, err)
				}
			})
			t.Run("Reader", func(t *testing.T) {
				_, err := Decode(errReader{}, 10).BigInt()
				require.Error(t, err)
			})
		})
	})
}

type errReader struct{}

func (errReader) Err() error {
	return io.ErrNoProgress
}

func (e errReader) Read([]byte) (int, error) {
	return 0, e.Err()
}
