package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// ObjectBytes calls f for every key in object, using byte slice as key.
//
// The key value is valid only until f is not returned.
func (d *Decoder) ObjectBytes(f func(d *Decoder, key []byte) error) error {
	if err := d.expectNext('{'); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := d.incrementDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := d.next()
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	if c == '}' {
		return d.decrementDepth()
	}
	d.unread()

	k, err := d.str(value{})
	if err != nil {
		return xerrors.Errorf("str: %w", err)
	}
	if err := d.expectNext(':'); err != nil {
		return xerrors.Errorf("field: %w", err)
	}
	if err := f(d, k.buf); err != nil {
		return xerrors.Errorf("callback: %w", err)
	}

	c, err = d.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	for c == ',' {
		k, err := d.str(value{})
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		if err := d.expectNext(':'); err != nil {
			return xerrors.Errorf("field: %w", err)
		}
		if err := f(d, k.buf); err != nil {
			return xerrors.Errorf("callback: %w", err)
		}
		if c, err = d.next(); err != nil {
			return xerrors.Errorf("next: %w", err)
		}
	}
	if c != '}' {
		return xerrors.Errorf("end: %w", badToken(c))
	}
	return d.decrementDepth()
}

// Object reads json object, calling f on each field.
//
// Use ObjectBytes to reduce heap allocations for keys.
func (d *Decoder) Object(f func(d *Decoder, key string) error) error {
	return d.ObjectBytes(func(d *Decoder, key []byte) error {
		return f(d, string(key))
	})
}
