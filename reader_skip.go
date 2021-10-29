package jx

import (
	"golang.org/x/xerrors"
)

// Null reads a json object as null and
// returns whether it's a null or not.
func (r *Reader) Null() error {
	if err := r.expectNext('n'); err != nil {
		return err
	}
	return r.skipThreeBytes('u', 'l', 'l') // null
}

// Bool reads a json object as Bool
func (r *Reader) Bool() (bool, error) {
	c, err := r.next()
	if err != nil {
		return false, err
	}
	switch c {
	case 't':
		if err := r.skipThreeBytes('r', 'u', 'e'); err != nil {
			return false, err
		}
		return true, nil
	case 'f':
		return false, r.skipFourBytes('a', 'l', 's', 'e')
	default:
		return false, badToken(c)
	}
}

// Skip skips a json object and positions to relatively the next json object.
func (r *Reader) Skip() error {
	c, err := r.next()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		if err := r.strSkip(); err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		return nil
	case 'n':
		return r.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		return r.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		return r.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		r.unread()
		_, err := r.Float32()
		return err
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return r.skipNumber()
	case '[':
		if err := r.skipArray(); err != nil {
			return xerrors.Errorf("array: %w", err)
		}
		return nil
	case '{':
		if err := r.skipObject(); err != nil {
			return xerrors.Errorf("object: %w", err)
		}
		return nil
	default:
		return badToken(c)
	}
}

func (r *Reader) skipFourBytes(b1, b2, b3, b4 byte) error {
	for _, b := range [...]byte{b1, b2, b3, b4} {
		c, err := r.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}

func (r *Reader) skipThreeBytes(b1, b2, b3 byte) error {
	for _, b := range [...]byte{b1, b2, b3} {
		c, err := r.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}
