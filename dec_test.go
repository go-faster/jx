package jx

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrNoProgress
}
