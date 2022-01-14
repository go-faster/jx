package jx

import (
	"github.com/go-faster/errors"
)

// ObjBytes calls f for every key in object, using byte slice as key.
//
// The key value is valid only until f is not returned.
func (d *Decoder) ObjBytes(f func(d *Decoder, key []byte) error) error {
	if err := d.consume('{'); err != nil {
		return errors.Wrap(err, "start")
	}
	if f == nil {
		return d.skipObj()
	}
	if err := d.incDepth(); err != nil {
		return errors.Wrap(err, "inc")
	}
	c, err := d.more()
	if err != nil {
		return errors.Wrap(err, "next")
	}
	if c == '}' {
		return d.decDepth()
	}
	d.unread()

	k, err := d.str(value{raw: true})
	if err != nil {
		return errors.Wrap(err, "str")
	}
	if err := d.consume(':'); err != nil {
		return errors.Wrap(err, "field")
	}
	// Skip whitespace.
	if _, err = d.more(); err != nil {
		return errors.Wrap(err, "more")
	}
	d.unread()
	if err := f(d, k.buf); err != nil {
		return errors.Wrap(err, "callback")
	}

	c, err = d.more()
	if err != nil {
		return errors.Wrap(err, "next")
	}
	for c == ',' {
		k, err := d.str(value{raw: true})
		if err != nil {
			return errors.Wrap(err, "str")
		}
		if err := d.consume(':'); err != nil {
			return errors.Wrap(err, "field")
		}
		// Check that value exists.
		if _, err = d.more(); err != nil {
			return errors.Wrap(err, "more")
		}
		d.unread()
		if err := f(d, k.buf); err != nil {
			return errors.Wrap(err, "callback")
		}
		if c, err = d.more(); err != nil {
			return errors.Wrap(err, "next")
		}
	}
	if c != '}' {
		return errors.Wrap(badToken(c), "err")
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
