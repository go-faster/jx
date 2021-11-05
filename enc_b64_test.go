package jx

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_Base64(t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		for _, s := range [][]byte{
			[]byte(`1`),
			[]byte(`12`),
			[]byte(`2345`),
			{1, 2, 3, 4, 5, 6},
		} {
			var e Encoder
			e.Base64(s)

			expected := fmt.Sprintf("%q", base64.StdEncoding.EncodeToString(s))
			require.Equal(t, expected, e.String())

			requireCompat(t, e.Bytes(), s)
		}
	})
	t.Run("Zero", func(t *testing.T) {
		var e Encoder
		e.Base64([]byte{})
		require.Equal(t, `null`, e.String())
	})
}

func BenchmarkEncoder_Base64(b *testing.B) {
	for _, n := range []int{
		128,
		256,
		512,
		1024,
	} {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			var v []byte
			for i := 0; i < n; i++ {
				v = append(v, byte(i%256))
			}

			b.ReportAllocs()
			b.SetBytes(int64(n))
			var e Encoder
			for i := 0; i < b.N; i++ {
				e.Base64(v)
				e.Reset()
			}
		})
	}
}
