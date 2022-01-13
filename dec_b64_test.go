package jx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Base64(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		for _, v := range [][]byte{
			[]byte("foo"),
			{1, 2, 3, 4},
			{1, 2, 3},
			{1, 2},
			{1},
			{},
			nil,
		} {
			var e Encoder
			e.Base64(v)

			d := DecodeBytes(e.Bytes())

			got, err := d.Base64()
			require.NoError(t, err)
			require.Equal(t, v, got)

			d.ResetBytes(e.Bytes())

			target := make([]byte, 0)
			if v == nil {
				// Append won't return nil, so just setting target
				// to nil to pass test.
				target = nil
			}
			got, err = d.Base64Append(target)
			require.NoError(t, err)
			require.Equal(t, v, got)
		}
	})
	t.Run("Negative", func(t *testing.T) {
		for _, v := range []string{
			`false`,
			`nu`,
			`12345`,
			`"foo`,
			`"100"`,
		} {
			t.Run(v, func(t *testing.T) {
				d := DecodeStr(v)

				_, err := d.Base64()
				require.Error(t, err)

				d = DecodeStr(v)
				_, err = d.Base64Append(nil)
				require.Error(t, err)
			})
		}
	})
}

func BenchmarkDecoder_Base64Append(b *testing.B) {
	for _, n := range []int{
		128,
		256,
		512,
		1024,
	} {
		b.Run(fmt.Sprintf("%db", n), func(b *testing.B) {
			var v []byte
			for i := 0; i < n; i++ {
				v = append(v, byte(i%256))
			}
			var e Encoder
			e.Base64(v)

			b.SetBytes(int64(n))
			b.ReportAllocs()

			d := DecodeBytes(nil)
			target := make([]byte, 0, n*2)
			for i := 0; i < b.N; i++ {
				d.ResetBytes(e.Bytes())
				if _, err := d.Base64Append(target); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
