package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// ObjBytes calls f for every key in object, using byte slice as key.
//
// The key value is valid only until f is not returned.
func (r *Reader) ObjBytes(f func(r *Reader, key []byte) error) error {
	if err := r.expectNext('{'); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := r.incrementDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := r.next()
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	if c == '}' {
		return r.decrementDepth()
	}
	r.unread()

	k, err := r.str(value{})
	if err != nil {
		return xerrors.Errorf("str: %w", err)
	}
	if err := r.expectNext(':'); err != nil {
		return xerrors.Errorf("field: %w", err)
	}
	if err := f(r, k.buf); err != nil {
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
		k, err := r.str(value{})
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		if err := r.expectNext(':'); err != nil {
			return xerrors.Errorf("field: %w", err)
		}
		if err := f(r, k.buf); err != nil {
			return xerrors.Errorf("callback: %w", err)
		}
		if c, err = r.next(); err != nil {
			return xerrors.Errorf("next: %w", err)
		}
	}
	if c != '}' {
		return xerrors.Errorf("end: %w", badToken(c))
	}
	return r.decrementDepth()
}

// Obj read json object, calling f on each field.
func (r *Reader) Obj(f func(i *Reader, key string) error) error {
	return r.ObjBytes(func(i *Reader, key []byte) error {
		return f(i, string(key))
	})
}
