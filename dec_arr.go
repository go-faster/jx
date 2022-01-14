package jx

import (
	"github.com/go-faster/errors"
)

// Elem skips to the start of next array element, returning true boolean
// if element exists.
//
// Can be called before or in Array.
func (d *Decoder) Elem() (ok bool, err error) {
	c, err := d.next()
	if err != nil {
		return false, err
	}
	switch c {
	case '[':
		c, err := d.more()
		if err != nil {
			return false, errors.Wrap(err, "next")
		}
		if c != ']' {
			d.unread()
			return true, nil
		}
		return false, nil
	case ']':
		return false, nil
	case ',':
		return true, nil
	default:
		return false, errors.Wrap(badToken(c), `"[" or "," or "]" expected`)
	}
}

// Arr decodes array and invokes callback on each array element.
func (d *Decoder) Arr(f func(d *Decoder) error) error {
	if err := d.consume('['); err != nil {
		return errors.Wrap(err, "start")
	}
	if f == nil {
		return d.skipArr()
	}
	if err := d.incDepth(); err != nil {
		return errors.Wrap(err, "inc")
	}
	c, err := d.more()
	if err != nil {
		return err
	}
	if c == ']' {
		return d.decDepth()
	}
	d.unread()
	if err := f(d); err != nil {
		return errors.Wrap(err, "callback")
	}

	c, err = d.more()
	if err != nil {
		return errors.Wrap(err, "next")
	}
	for c == ',' {
		// Skip whitespace before reading element.
		if _, err := d.next(); err != nil {
			return errors.Wrap(err, "next")
		}
		d.unread()
		if err := f(d); err != nil {
			return errors.Wrap(err, "callback")
		}
		if c, err = d.next(); err != nil {
			return errors.Wrap(err, "next")
		}
	}
	if c != ']' {
		return errors.Wrap(badToken(c), "end")
	}
	return d.decDepth()
}
