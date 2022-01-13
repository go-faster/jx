package jx

import (
	"io"

	"github.com/go-faster/errors"
)

// Null reads a json object as null and
// returns whether it's a null or not.
func (d *Decoder) Null() error {
	if err := d.consume('n'); err != nil {
		return err
	}
	return d.skipThreeBytes('u', 'l', 'l') // null
}

// Bool reads a json object as Bool
func (d *Decoder) Bool() (bool, error) {
	c, err := d.next()
	if err != nil {
		return false, err
	}
	switch c {
	case 't':
		if err := d.skipThreeBytes('r', 'u', 'e'); err != nil {
			return false, err
		}
		return true, nil
	case 'f':
		return false, d.skipFourBytes('a', 'l', 's', 'e')
	default:
		return false, badToken(c)
	}
}

// Skip skips a json object and positions to relatively the next json object.
func (d *Decoder) Skip() error {
	c, err := d.next()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		if err := d.skipStr(); err != nil {
			return errors.Wrap(err, "str")
		}
		return nil
	case 'n':
		return d.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		return d.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		return d.skipFourBytes('a', 'l', 's', 'e') // false
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		d.unread()
		return d.skipNumber()
	case '[':
		if err := d.skipArr(); err != nil {
			return errors.Wrap(err, "array")
		}
		return nil
	case '{':
		if err := d.skipObj(); err != nil {
			return errors.Wrap(err, "object")
		}
		return nil
	default:
		return badToken(c)
	}
}

func (d *Decoder) skipFourBytes(b1, b2, b3, b4 byte) error {
	for _, b := range [...]byte{b1, b2, b3, b4} {
		c, err := d.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}

func (d *Decoder) skipThreeBytes(b1, b2, b3 byte) error {
	for _, b := range [...]byte{b1, b2, b3} {
		c, err := d.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}

var (
	closerSet = [256]byte{
		',':  1,
		']':  1,
		'}':  1,
		' ':  1,
		'\t': 1,
		'\n': 1,
		'\r': 1,
	}
	digitSet = [256]byte{
		'0': 1,
		'1': 1,
		'2': 1,
		'3': 1,
		'4': 1,
		'5': 1,
		'6': 1,
		'7': 1,
		'8': 1,
		'9': 1,
	}
)

// skipNumber reads one JSON number.
//
// Assumes d.buf is not empty.
func (d *Decoder) skipNumber() error {
	c := d.buf[d.head]
	d.head++
	switch c {
	case '-':
		c, err := d.byte()
		if err != nil {
			return err
		}
		// Character after '-' must be a digit.
		if digitSet[c] == 0 {
			return badToken(c)
		}
	case '0':
		// If buffer is empty, try to read more.
		if d.head == d.tail {
			err := d.read()
			if err != nil {
				// There is no data anymore.
				if err == io.EOF {
					return nil
				}
				return err
			}
		}

		c = d.buf[d.head]
		if closerSet[c] != 0 {
			return nil
		}
		switch c {
		case '.':
			goto stateDot
		case 'e', 'E':
			goto stateExp
		default:
			return badToken(c)
		}
	}
	for {
		for i, c := range d.buf[d.head:d.tail] {
			if closerSet[c] != 0 {
				d.head += i
				return nil
			}
			if digitSet[c] != 0 {
				continue
			}
			switch c {
			case '.':
				d.head += i
				goto stateDot
			case 'e', 'E':
				d.head += i
				goto stateExp
			default:
				return badToken(c)
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

stateDot:
	d.head++
	for {
		var last byte = '.'
		for i, c := range d.buf[d.head:d.tail] {
			if closerSet[c] != 0 {
				d.head += i
				if last == '.' {
					return io.ErrUnexpectedEOF
				}
				return nil
			}
			last = c
			if digitSet[c] != 0 {
				continue
			}
			switch c {
			case 'e', 'E':
				if i == 0 {
					return badToken(c)
				}
				d.head += i
				goto stateExp
			default:
				return badToken(c)
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				if last == '.' {
					return io.ErrUnexpectedEOF
				}
				return nil
			}
			return err
		}
	}
stateExp:
	d.head++
	// There must be a number or sign after e.
	{
		v, err := d.byte()
		if err != nil {
			return err
		}
		if digitSet[v] == 0 {
			// There must be a number after e.
			if v == '-' || v == '+' {
				v, err := d.byte()
				if err != nil {
					return err
				}
				if digitSet[v] == 0 {
					return badToken(v)
				}
			} else {
				return badToken(v)
			}
		}
	}
	for {
		for i, c := range d.buf[d.head:d.tail] {
			if closerSet[c] != 0 {
				d.head += i
				return nil
			}
			if digitSet[c] == 0 {
				return badToken(c)
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (d *Decoder) skipStr() error {
	ok, err := d.skipStrFast()
	if err != nil || ok {
		return err
	}

	d.unread()
	if _, err := d.str(value{ignore: true}); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) skipStrFast() (ok bool, err error) {
	for i := d.head; i < d.tail; i++ {
		c := d.buf[i]
		switch {
		case c == '"':
			d.head = i + 1
			return true, nil
		case c == '\\':
			return false, nil
		case c < ' ':
			return false, badToken(c)
		}
	}
	return false, nil
}

func (d *Decoder) skipObj() error {
	d.unread()
	return d.Obj(nil)
}

func (d *Decoder) skipArr() error {
	d.unread()
	return d.Arr(nil)
}
