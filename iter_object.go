package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// ObjBytes calls f for every key in object, using byte slice as key.
//
// The key value is valid only until f is not returned.
func (it *Iter) ObjBytes(f func(i *Iter, key []byte) error) error {
	if err := it.expectNext('{'); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := it.incrementDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := it.next()
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	if c == '}' {
		return it.decrementDepth()
	}
	it.unread()

	k, err := it.str(value{})
	if err != nil {
		return xerrors.Errorf("str: %w", err)
	}
	if err := it.expectNext(':'); err != nil {
		return xerrors.Errorf("field: %w", err)
	}
	if err := f(it, k.buf); err != nil {
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
		k, err := it.str(value{})
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		if err := it.expectNext(':'); err != nil {
			return xerrors.Errorf("field: %w", err)
		}
		if err := f(it, k.buf); err != nil {
			return xerrors.Errorf("callback: %w", err)
		}
		if c, err = it.next(); err != nil {
			return xerrors.Errorf("next: %w", err)
		}
	}
	if c != '}' {
		return xerrors.Errorf("end: %w", badToken(c))
	}
	return it.decrementDepth()
}

// Obj read json object, calling f on each field.
func (it *Iter) Obj(f func(i *Iter, key string) error) error {
	return it.ObjBytes(func(i *Iter, key []byte) error {
		return f(i, string(key))
	})
}
