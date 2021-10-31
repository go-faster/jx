//go:build gofuzz
// +build gofuzz

package jx

import (
	"bytes"
	"fmt"
)

func Fuzz(data []byte) int {
	_ = Valid(data)
	return 1
}

var (
	d = GetDecoder()
	w = GetEncoder()
)

func FuzzAny(data []byte) int {
	d.ResetBytes(data)

	v, err := d.Any()
	if err != nil {
		return 0
	}
	if v.Type == AnyInvalid {
		panic(v.Type)
	}

	w.Reset()
	w.Any(v)

	// Parsing from buf to new value.
	d.ResetBytes(w.Bytes())
	parsed, err := d.Any()
	if err != nil {
		fmt.Printf("%v:\nBuf:   %s\nValue: %s\nData:  %s",
			err, w.Bytes(), v, data)
		panic(err)
	}
	if !parsed.Equal(v) {
		fmt.Printf("\nBuf:   %s\nValue: %s != %s \nData:  %s",
			w.Bytes(), parsed, v, data)
		panic("not equal")
	}
	b := w.Bytes()
	w.SetBytes(nil)
	parsed.Write(w)
	if !bytes.Equal(w.Bytes(), b) {
		panic(fmt.Sprintf("%s != %s", w, b))
	}

	return 1
}
