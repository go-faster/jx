package jx

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncoder_Base64(t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		for i, s := range [][]byte{
			[]byte(`1`),
			[]byte(`12`),
			[]byte(`2345`),
			{1, 2, 3, 4, 5, 6},

			bytes.Repeat([]byte{1}, encoderBufSize-1),
			bytes.Repeat([]byte{1}, encoderBufSize),
			bytes.Repeat([]byte{1}, encoderBufSize+1),
		} {
			s := s
			t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
				requireCompat(t, func(e *Encoder) {
					e.Base64(s)
				}, s)
			})
		}
	})
	t.Run("Zeroes", func(t *testing.T) {
		t.Run("Nil", func(t *testing.T) {
			s := []byte(nil)
			requireCompat(t, func(e *Encoder) {
				e.Base64(s)
			}, s)
		})
		t.Run("ZeroLen", func(t *testing.T) {
			s := make([]byte, 0)
			requireCompat(t, func(e *Encoder) {
				e.Base64(s)
			}, s)
		})
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
