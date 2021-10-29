package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// Elem reads array element and reports whether array has more
// elements to read.
func (it *Iter) Elem() (ok bool, err error) {
	c, err := it.next()
	if err != nil {
		return false, err
	}
	switch c {
	case '[':
		c, err := it.next()
		if err != nil {
			return false, xerrors.Errorf("next: %w", err)
		}
		if c != ']' {
			it.unread()
			return true, nil
		}
		return false, nil
	case ']':
		return false, nil
	case ',':
		return true, nil
	default:
		return false, xerrors.Errorf(`"[" or "," or "]" expected: %w`, badToken(c))
	}
}

// Array reads array and call f on each element.
func (it *Iter) Array(f func(i *Iter) error) error {
	if err := it.expectNext('['); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := it.incrementDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := it.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return err
	}
	if c == ']' {
		return it.decrementDepth()
	}
	it.unread()
	if err := f(it); err != nil {
		return xerrors.Errorf("callback: %w", err)
	}

	c, err = it.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	for c == ',' {
		if err := f(it); err != nil {
			return xerrors.Errorf("callback: %w", err)
		}
		if c, err = it.next(); err != nil {
			return xerrors.Errorf("next: %w", err)
		}
	}
	if c != ']' {
		return xerrors.Errorf("end: %w", badToken(c))
	}
	return it.decrementDepth()
}
