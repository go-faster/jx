package byteseq

import (
	"testing"
	"unicode/utf8"
)

func BenchmarkDecodeRuneInByteseq(b *testing.B) {
	var (
		buf    [4]byte
		result rune
	)
	utf8.EncodeRune(buf[:], 'ж')

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, _ = DecodeRuneInByteseq(buf[:])
	}

	if result != 'ж' {
		b.Fatal(result)
	}
}
