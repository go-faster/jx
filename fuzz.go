//go:build gofuzz
// +build gofuzz

package jx

func Fuzz(data []byte) int {
	_ = Valid(data)
	return 1
}
