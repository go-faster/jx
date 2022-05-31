//go:build !go1.18

package jx

import "unicode/utf8"

func decodeRuneInByteseq(val string) (r rune, size int) {
	var tmp [4]byte
	n := copy(tmp[:], val)
	return utf8.DecodeRune(tmp[:n])
}
