//go:build gofuzz
// +build gofuzz

package jx

import (
	"bytes"
	"fmt"
	"reflect"
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
	if err := w.Any(v); err != nil {
		panic(errMaxDepth)
	}

	// Parsing from buf to new value.
	d.ResetBytes(w.Bytes())
	parsed, err := d.Any()
	if err != nil {
		fmt.Printf("%v:\nBuf:   %s\nValue: %s\nData:  %s",
			err, w.Bytes(), v, data)
		panic(err)
	}
	if !reflect.DeepEqual(parsed, v) {
		fmt.Printf("%v:\nBuf:   %s\nValue: %s != %s \nData:  %s",
			nil, w.Bytes(), parsed, v, data)
		panic("not equal")
	}
	b := w.Bytes()
	w.SetBytes(nil)
	if err := parsed.Write(w); err != nil {
		panic(err)
	}
	if !bytes.Equal(w.Bytes(), b) {
		panic(fmt.Sprintf("%s != %s", w, b))
	}

	return 1
}
