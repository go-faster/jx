package jx

import (
	"testing"
	"unicode/utf8"
)

func Benchmark_decodeRuneInByteseq(b *testing.B) {
	var (
		buf    [4]byte
		result rune
	)
	utf8.EncodeRune(buf[:], 'ж')

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, _ = decodeRuneInByteseq(buf[:])
	}

	if result != 'ж' {
		b.Fatal(result)
	}
}
