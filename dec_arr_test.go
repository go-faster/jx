package jx

import (
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Arr(t *testing.T) {
	t.Run("Blank", func(t *testing.T) {
		d := DecodeStr(`[]`)
		require.NoError(t, d.Arr(nil))
	})
	t.Run("Invalid", func(t *testing.T) {
		d := DecodeStr(`{`)
		require.Error(t, d.Arr(nil))
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("")
		require.ErrorIs(t, d.Arr(nil), io.ErrUnexpectedEOF)
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("[")
		require.ErrorIs(t, d.Arr(nil), io.ErrUnexpectedEOF)
	})
	t.Run("Invalid", func(t *testing.T) {
		for _, s := range testArrs {
			checker := require.Error
			if json.Valid([]byte(s)) {
				checker = require.NoError
			}

			d := DecodeStr(s)
			checker(t, d.Arr(crawlValue), s)
		}
	})
	t.Run("Whitespace", func(t *testing.T) {
		d := DecodeStr(`[1 , 2,  3 ,45, 6]`)
		require.NoError(t, d.Arr(func(d *Decoder) error {
			_, err := d.Int()
			return err
		}))
	})
	t.Run("Depth", func(t *testing.T) {
		var data []byte
		for i := 0; i <= maxDepth; i++ {
			data = append(data, '[')
		}
		d := DecodeBytes(data)
		require.ErrorIs(t, d.Arr(nil), errMaxDepth)
	})
}

func TestDecoder_Elem(t *testing.T) {
	t.Run("Blank", func(t *testing.T) {
		d := DecodeStr(`[]`)
		ok, err := d.Elem()
		require.NoError(t, err)
		require.False(t, ok)
	})
	t.Run("Invalid", func(t *testing.T) {
		d := DecodeStr(`{`)
		ok, err := d.Elem()
		require.Error(t, err)
		require.False(t, ok)
	})
	t.Run("EOF", func(t *testing.T) {
		d := DecodeStr("")
		ok, err := d.Elem()
		require.ErrorIs(t, err, io.EOF)
		require.False(t, ok)
	})
	t.Run("ErrUnexpectedEOF", func(t *testing.T) {
		d := DecodeStr("[")
		ok, err := d.Elem()
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.False(t, ok)
	})
}

//go:embed testdata/bools.json
var boolsData []byte

func BenchmarkDecodeBools(b *testing.B) {
	b.Run("Callback", func(b *testing.B) {
		d := DecodeBytes(boolsData)
		r := make([]bool, 0, 100)

		b.SetBytes(int64(len(boolsData)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r = r[:0]
			d.ResetBytes(boolsData)

			if err := d.Arr(func(d *Decoder) error {
				f, err := d.Bool()
				if err != nil {
					return err
				}
				r = append(r, f)
				return nil
			}); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Iterator", func(b *testing.B) {
		d := DecodeBytes(boolsData)
		r := make([]bool, 0, 100)

		b.SetBytes(int64(len(boolsData)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r = r[:0]
			d.ResetBytes(boolsData)

			iter, err := d.ArrIter()
			if err != nil {
				b.Fatal(err)
			}
			for iter.Next() {
				v, err := d.Bool()
				if err != nil {
					b.Fatal(err)
				}
				r = append(r, v)
			}
			if err := iter.Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}
