package jx

import (
	"github.com/go-faster/errors"
)

// Capture calls f and then rolls back to state before call.
//
// Does not work with reader.
func (d *Decoder) Capture(f func(d *Decoder) error) error {
	if d.reader != nil {
		return errors.New("capture is not supported with reader")
	}
	if f == nil {
		return nil
	}
	head, tail, depth := d.head, d.tail, d.depth
	err := f(d)
	d.head, d.tail, d.depth = head, tail, depth
	return err
}
