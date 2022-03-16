//go:build appengine || purego

package jx

import (
	"strconv"

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

	val, err := strconv.ParseFloat(string(str), size)
	if err != nil {
		return 0, err
	}

	return val, nil
}
