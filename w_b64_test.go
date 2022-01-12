package jx

import (
	"fmt"
	"testing"
)

func BenchmarkWriter_Base64(b *testing.B) {
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
			var w Writer
			for i := 0; i < b.N; i++ {
				w.Base64(v)
				w.Reset()
			}
		})
	}
}
