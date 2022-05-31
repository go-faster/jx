//go:build !go1.18

package jx

import (
	"testing"
)

func Benchmark_decodeRuneInByteseq(b *testing.B) {
	var result rune
	const buf = `ж`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, _ = decodeRuneInByteseq(buf[:])
	}

	if result != 'ж' {
		b.Fatal(result)
	}
}
