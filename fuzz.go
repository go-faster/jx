//go:build gofuzz
// +build gofuzz

package jx

import (
	"encoding/json"
	"fmt"
)

func Fuzz(data []byte) int {
	got := Valid(data)
	exp := json.Valid(data)
	if !exp && got {
		fmt.Printf("jx: %v\nencoding/json:%v\n", got, exp)
		panic("mismatch")
	}
	return 1
}
