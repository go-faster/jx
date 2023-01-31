package byteseq

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestDecodeRuneInByteseq(t *testing.T) {
	for i, tt := range []struct {
		s string
	}{
		{""},
		{"\x00"},
		{"a"},
		{"ж"},
		{"🤡"},
		{"👩🏿"},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			gotR, gotSize := DecodeRuneInByteseq(tt.s)
			expectR, expectSize := utf8.DecodeRuneInString(tt.s)

			require.Equal(t, expectR, gotR)
			require.Equal(t, expectSize, gotSize)
		})
	}
}

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
