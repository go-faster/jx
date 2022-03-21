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

func createTestCase(input string, cb func(t *testing.T, d *Decoder) error) func(t *testing.T) {
	run := func(d *Decoder, input string, valid bool) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	}

	return func(t *testing.T) {
		valid := json.Valid([]byte(input))

		t.Run("Buffer", run(DecodeStr(input), input, valid))

		r := strings.NewReader(input)
		d := Decode(r, 512)
		t.Run("Reader", run(d, input, valid))

		r.Reset(input)
		obr := iotest.OneByteReader(r)
		d.Reset(obr)
		t.Run("OneByteReader", run(d, input, valid))
	}
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

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, e.Err()
}
