package jx

import (
	"golang.org/x/xerrors"
)

// Capture calls f and then rolls back buffer to state before call.
//
// Does not work with reader.
func (it *Iter) Capture(f func(i *Iter) error) error {
	if it.reader != nil {
		return xerrors.New("capture is not supported with reader")
	}
	head, tail, depth := it.head, it.tail, it.depth
	err := f(it)
	it.head, it.tail, it.depth = head, tail, depth
	return err
}
