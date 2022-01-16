package jx

import (
	"io"
	"math/bits"

	"github.com/go-faster/errors"
)

func (d *Decoder) readExact4(b *[4]byte) error {
	if buf := d.buf[d.head:d.tail]; len(buf) >= len(b) {
		d.head += copy(b[:], buf[:4])
		return nil
	}

	n := copy(b[:], d.buf[d.head:d.tail])
	if err := d.readAtLeast(len(b) - n); err != nil {
		return err
	}
	d.head += copy(b[n:], d.buf[d.head:d.tail])
	return nil
}

func findInvalidToken4(buf [4]byte, mask uint32) error {
	c := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
	idx := bits.TrailingZeros32(c^mask) / 8
	return badToken(buf[idx])
}

// Null reads a json object as null and
// returns whether it's a null or not.
func (d *Decoder) Null() error {
	var buf [4]byte
	if err := d.readExact4(&buf); err != nil {
		return err
	}

	if string(buf[:]) != "null" {
		const encodedNull = 'n' | 'u'<<8 | 'l'<<16 | 'l'<<24
		return findInvalidToken4(buf, encodedNull)
	}
	return nil
}

// Bool reads a json object as Bool
func (d *Decoder) Bool() (bool, error) {
	var buf [4]byte
	if err := d.readExact4(&buf); err != nil {
		return false, err
	}

	switch string(buf[:]) {
	case "true":
		return true, nil
	case "fals":
		c, err := d.byte()
		if err != nil {
			return false, err
		}
		if c != 'e' {
			return false, badToken(c)
		}
		return false, nil
	default:
		switch c := buf[0]; c {
		case 't':
			const encodedTrue = 't' | 'r'<<8 | 'u'<<16 | 'e'<<24
			return false, findInvalidToken4(buf, encodedTrue)
		case 'f':
			const encodedAlse = 'a' | 'l'<<8 | 's'<<16 | 'e'<<24
			return false, findInvalidToken4(buf, encodedAlse)
		default:
			return false, badToken(c)
		}
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
		d.unread()
		return d.Null()
	case 't', 'f':
		d.unread()
		_, err := d.Bool()
		return err
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
		'E': 1, 'F': 1,

		'a': 1, 'b': 1, 'c': 1, 'd': 1,
		'e': 1, 'f': 1,
	}
)

// skipStr reads one JSON string.
//
// Assumes first quote was consumed.
func (d *Decoder) skipStr() error {
	var (
		c byte
		i int
	)
readStr:
	for {
		i = 0
		buf := d.buf[d.head:d.tail]
		for len(buf) >= 8 {
			c = buf[0]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[1]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[2]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[3]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[4]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[5]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[6]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			c = buf[7]
			if safeSet[c] != 0 {
				goto readTok
			}
			i++

			buf = buf[8:]
		}
		var n int
		for n, c = range buf {
			if safeSet[c] != 0 {
				i += n
				goto readTok
			}
		}

		if err := d.read(); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return err
		}
	}

readTok:
	switch {
	case c == '"':
		d.head += i + 1
		return nil
	case c == '\\':
		d.head += i + 1
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
	case c < ' ':
		return badToken(c)
	}
	goto readStr
}

// skipObj reads JSON object.
//
// Assumes first bracket was consumed.
func (d *Decoder) skipObj() error {
	if err := d.incDepth(); err != nil {
		return errors.Wrap(err, "inc")
	}

	c, err := d.more()
	if err != nil {
		return errors.Wrap(err, "next")
	}
	switch c {
	case '}':
		return d.decDepth()
	case '"':
		d.unread()
	default:
		return badToken(c)
	}

	for {
		if err := d.consume('"'); err != nil {
			return err
		}
		if err := d.skipStr(); err != nil {
			return errors.Wrap(err, "read field name")
		}
		if err := d.consume(':'); err != nil {
			return errors.Wrap(err, "field")
		}
		if err := d.Skip(); err != nil {
			return err
		}
		c, err := d.more()
		if err != nil {
			return errors.Wrap(err, "read comma")
		}
		switch c {
		case ',':
			continue
		case '}':
			return d.decDepth()
		default:
			return badToken(c)
		}
	}
}

// skipArr reads JSON array.
//
// Assumes first bracket was consumed.
func (d *Decoder) skipArr() error {
	if err := d.incDepth(); err != nil {
		return errors.Wrap(err, "inc")
	}

	c, err := d.more()
	if err != nil {
		return errors.Wrap(err, "next")
	}
	if c == ']' {
		return d.decDepth()
	}
	d.unread()

	for {
		if err := d.Skip(); err != nil {
			return err
		}
		c, err := d.more()
		if err != nil {
			return errors.Wrap(err, "read comma")
		}
		switch c {
		case ',':
			continue
		case ']':
			return d.decDepth()
		default:
			return badToken(c)
		}
	}
}
