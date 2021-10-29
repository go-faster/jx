package jx

import (
	"golang.org/x/xerrors"
)

// Null reads a json object as null and
// returns whether it's a null or not.
func (it *Iter) Null() error {
	if err := it.expectNext('n'); err != nil {
		return err
	}
	return it.skipThreeBytes('u', 'l', 'l') // null
}

// Bool reads a json object as Bool
func (it *Iter) Bool() (bool, error) {
	c, err := it.next()
	if err != nil {
		return false, err
	}
	switch c {
	case 't':
		if err := it.skipThreeBytes('r', 'u', 'e'); err != nil {
			return false, err
		}
		return true, nil
	case 'f':
		return false, it.skipFourBytes('a', 'l', 's', 'e')
	default:
		return false, badToken(c)
	}
}

// Skip skips a json object and positions to relatively the next json object.
func (it *Iter) Skip() error {
	c, err := it.next()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		if err := it.strSkip(); err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		return nil
	case 'n':
		return it.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		return it.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		return it.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		it.unread()
		_, err := it.Float32()
		return err
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return it.skipNumber()
	case '[':
		if err := it.skipArray(); err != nil {
			return xerrors.Errorf("array: %w", err)
		}
		return nil
	case '{':
		if err := it.skipObject(); err != nil {
			return xerrors.Errorf("object: %w", err)
		}
		return nil
	default:
		return badToken(c)
	}
}

func (it *Iter) skipFourBytes(b1, b2, b3, b4 byte) error {
	if err := it.skipThreeBytes(b1, b2, b3); err != nil {
		return err
	}
	if it.byte() != b4 {
		return badToken(it.byte())
	}
	return nil
}

func (it *Iter) skipThreeBytes(b1, b2, b3 byte) error {
	if it.byte() != b1 {
		return badToken(it.byte())
	}
	if it.byte() != b2 {
		return badToken(it.byte())
	}
	if it.byte() != b3 {
		return badToken(it.byte())
	}
	return nil
}
