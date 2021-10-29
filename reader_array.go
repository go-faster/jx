package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// Elem reads array element and reports whether array has more
// elements to read.
func (r *Reader) Elem() (ok bool, err error) {
	c, err := r.next()
	if err != nil {
		return false, err
	}
	switch c {
	case '[':
		c, err := r.next()
		if err != nil {
			return false, xerrors.Errorf("next: %w", err)
		}
		if c != ']' {
			r.unread()
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
func (r *Reader) Array(f func(r *Reader) error) error {
	if err := r.expectNext('['); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := r.incrementDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := r.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return err
	}
	if c == ']' {
		return r.decrementDepth()
	}
	r.unread()
	if err := f(r); err != nil {
		return xerrors.Errorf("callback: %w", err)
	}

	c, err = r.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	for c == ',' {
		if err := f(r); err != nil {
			return xerrors.Errorf("callback: %w", err)
		}
		if c, err = r.next(); err != nil {
			return xerrors.Errorf("next: %w", err)
		}
	}
	if c != ']' {
		return xerrors.Errorf("end: %w", badToken(c))
	}
	return r.decrementDepth()
}
