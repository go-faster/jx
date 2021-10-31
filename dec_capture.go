package jx

import (
	"golang.org/x/xerrors"
)

// Capture calls f and then rolls back buffer to state before call.
//
// Does not work with reader.
func (d *Decoder) Capture(f func(i *Decoder) error) error {
	if d.reader != nil {
		return xerrors.New("capture is not supported with reader")
	}
	if f == nil {
		return nil
	}
	head, tail, depth := d.head, d.tail, d.depth
	err := f(d)
	d.head, d.tail, d.depth = head, tail, depth
	return err
}
