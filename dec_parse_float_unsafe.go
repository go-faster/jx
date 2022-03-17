//go:build !appengine && !purego

package jx

import (
	"strconv"
	"unsafe"

	"github.com/go-faster/errors"
)

func (d *Decoder) floatSlow(size int) (float64, error) {
	var buf [32]byte

	str, err := d.numberAppend(buf[:0])
	if err != nil {
		return 0, errors.Wrap(err, "number")
	}
	if err := validateFloat(str); err != nil {
		return 0, errors.Wrap(err, "invalid")
	}

	slice := *(*sliceType)(unsafe.Pointer(&str)) // #nosec G103
	s := strType{
		Ptr: noescape(slice.Ptr),
		Len: slice.Len,
	}
	val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&s)), size) // #nosec G103
	if err != nil {
		return 0, err
	}

	return val, nil
}
