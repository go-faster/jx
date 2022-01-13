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
	skipNumberSet = [256]byte{
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

		',':  2,
		']':  2,
		'}':  2,
		' ':  2,
		'\t': 2,
		'\n': 2,
		'\r': 2,
	}
)

// skipNumber reads one JSON number.
//
// Assumes d.buf is not empty.
func (d *Decoder) skipNumber() error {
	const (
		digitTag  byte = 1
		closerTag byte = 2
	)
	c := d.buf[d.head]
	d.head++
	switch c {
	case '-':
		c, err := d.byte()
		if err != nil {
			return err
		}
		// Character after '-' must be a digit.
		if skipNumberSet[c] != digitTag {
			return badToken(c)
		}
		if c != '0' {
			break
		}
		fallthrough
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
		if skipNumberSet[c] == closerTag {
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
			switch skipNumberSet[c] {
			case closerTag:
				d.head += i
				return nil
			case digitTag:
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
				d.head = d.tail
				return nil
			}
			return err
		}
	}

stateDot:
	d.head++
	{
		var last byte = '.'
		for {
			for i, c := range d.buf[d.head:d.tail] {
				switch skipNumberSet[c] {
				case closerTag:
					d.head += i
					// Check that dot is not last character.
					if last == '.' {
						return io.ErrUnexpectedEOF
					}
					return nil
				case digitTag:
					last = c
					continue
				}

				switch c {
				case 'e', 'E':
					if last == '.' {
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
					d.head = d.tail
					// Check that dot is not last character.
					if last == '.' {
						return io.ErrUnexpectedEOF
					}
					return nil
				}
				return err
			}
		}
	}
stateExp:
	d.head++
	// There must be a number or sign after e.
	{
		numOrSign, err := d.byte()
		if err != nil {
			return err
		}
		if skipNumberSet[numOrSign] != digitTag { // If next character is not a digit, check for sign.
			if numOrSign == '-' || numOrSign == '+' {
				num, err := d.byte()
				if err != nil {
					return err
				}
				// There must be a number after sign.
				if skipNumberSet[num] != digitTag {
					return badToken(num)
				}
			} else {
				return badToken(numOrSign)
			}
		}
	}
	for {
		for i, c := range d.buf[d.head:d.tail] {
			if skipNumberSet[c] == closerTag {
				d.head += i
				return nil
			}
			if skipNumberSet[c] == 0 {
				return badToken(c)
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				d.head = d.tail
				return nil
			}
			return err
		}
	}
}

var (
	escapedStrSet = [256]byte{
		'"': 1, '\\': 1, '/': 1, 'b': 1, 'f': 1, 'n': 1, 'r': 1, 't': 1,
		'u': 2,
	}
	hexSet = [256]byte{
		'0': 1, '1': 1, '2': 1, '3': 1,
		'4': 1, '5': 1, '6': 1, '7': 1,
		'8': 1, '9': 1,

		'A': 1, 'B': 1, 'C': 1, 'D': 1,
		'E': 1, 'F': 1, 'G': 1, 'H': 1,
		'I': 1, 'J': 1, 'K': 1, 'L': 1,
		'M': 1, 'N': 1, 'O': 1, 'P': 1,
		'Q': 1, 'R': 1, 'S': 1, 'T': 1,
		'U': 1, 'V': 1, 'W': 1, 'X': 1,
		'Y': 1, 'Z': 1,

		'a': 1, 'b': 1, 'c': 1, 'd': 1,
		'e': 1, 'f': 1, 'g': 1, 'h': 1,
		'i': 1, 'j': 1, 'k': 1, 'l': 1,
		'm': 1, 'n': 1, 'o': 1, 'p': 1,
		'q': 1, 'r': 1, 's': 1, 't': 1,
		'u': 1, 'v': 1, 'w': 1, 'x': 1,
		'y': 1, 'z': 1,
	}
)

func (d *Decoder) skipStr() error {
readStr:
	for {
		for i, c := range d.buf[d.head:d.tail] {
			switch {
			case c == '"':
				d.head += i + 1
				return nil
			case c == '\\':
				d.head += i + 1
				goto readEscaped
			case c < ' ':
				return badToken(c)
			}
		}

		if err := d.read(); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return err
		}
	}

readEscaped:
	v, err := d.byte()
	if err != nil {
		return err
	}
	switch escapedStrSet[v] {
	case 1:
	case 2:
		for i := 0; i < 4; i++ {
			h, err := d.byte()
			if err != nil {
				return err
			}
			if hexSet[h] == 0 {
				return badToken(h)
			}
		}
	default:
		return badToken(v)
	}
	goto readStr
}

func (d *Decoder) skipObj() error {
	d.unread()
	return d.Obj(nil)
}

func (d *Decoder) skipArr() error {
	d.unread()
	return d.Arr(nil)
}
