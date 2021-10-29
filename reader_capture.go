package jx

import (
	"golang.org/x/xerrors"
)

// Capture calls f and then rolls back buffer to state before call.
//
// Does not work with reader.
func (r *Reader) Capture(f func(i *Reader) error) error {
	if r.reader != nil {
		return xerrors.New("capture is not supported with reader")
	}
	head, tail, depth := r.head, r.tail, r.depth
	err := f(r)
	r.head, r.tail, r.depth = head, tail, depth
	return err
}
