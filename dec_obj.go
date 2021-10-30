package jx

import (
	"io"

	"golang.org/x/xerrors"
)

func skipObjBytes(d *Decoder, _ []byte) error { return d.Skip() }

// ObjBytes calls f for every key in object, using byte slice as key.
//
// The key value is valid only until f is not returned.
func (d *Decoder) ObjBytes(f func(d *Decoder, key []byte) error) error {
	skip := f == nil
	if skip {
		// No callback is provided, skipping object.
		// Drop keys, skip values.
		f = skipObjBytes
	}
	if err := d.consume('{'); err != nil {
		return xerrors.Errorf("start: %w", err)
	}
	if err := d.incDepth(); err != nil {
		return xerrors.Errorf("inc: %w", err)
	}
	c, err := d.next()
	if err != nil {
		return xerrors.Errorf("next: %w", err)
	}
	if c == '}' {
		return d.decDepth()
	}
	d.unread()

	k, err := d.str(value{ignore: skip})
	if err != nil {
		return xerrors.Errorf("str: %w", err)
	}
	if err := d.consume(':'); err != nil {
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
		k, err := d.str(value{ignore: skip})
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		if err := d.consume(':'); err != nil {
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
	return d.decDepth()
}

// Obj reads json object, calling f on each field.
//
// Use ObjBytes to reduce heap allocations for keys.
func (d *Decoder) Obj(f func(d *Decoder, key string) error) error {
	if f == nil {
		// Skipping object.
		return d.ObjBytes(nil)
	}
	return d.ObjBytes(func(d *Decoder, key []byte) error {
		return f(d, string(key))
	})
}
