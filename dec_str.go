package jx

import (
	"fmt"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/go-faster/errors"
)

// StrAppend reads string and appends it to byte slice.
func (d *Decoder) StrAppend(b []byte) ([]byte, error) {
	v := value{
		buf: b,
		raw: false,
	}
	var err error
	if v, err = d.str(v); err != nil {
		return b, err
	}
	return v.buf, nil
}

type value struct {
	buf    []byte
	raw    bool // false forces buf reuse
	ignore bool
}

func (v value) rune(r rune) value {
	if v.ignore {
		return v
	}
	return value{
		buf: appendRune(v.buf, r),
		raw: v.raw,
	}
}

func (v value) byte(b byte) value {
	if v.ignore {
		return v
	}
	return value{
		buf: append(v.buf, b),
		raw: v.raw,
	}
}

// badTokenErr means that Token was unexpected while decoding.
type badTokenErr struct {
	Token byte
}

func (e badTokenErr) Error() string {
	return fmt.Sprintf("unexpected byte %d '%s'", e.Token, []byte{e.Token})
}

func badToken(c byte) error {
	return badTokenErr{Token: c}
}

func (d *Decoder) str(v value) (value, error) {
	if err := d.consume('"'); err != nil {
		return value{}, errors.Wrap(err, "start")
	}
	buf := d.buf[d.head:d.tail]
	for i, c := range buf {
		if c == '\\' {
			// Character is escaped, fallback to slow path.
			break
		}
		if c == '"' {
			// End of string in fast path.
			if v.ignore {
				d.head += i + 1
				return value{}, nil
			}
			str := buf[:i]
			d.head += i + 1
			if v.raw {
				return value{buf: str}, nil
			}
			return value{buf: append(v.buf, str...)}, nil
		}
		if c < ' ' {
			return value{}, errors.Wrap(badToken(c), "control character")
		}
	}
	return d.strSlow(v)
}

// StrBytes returns string value as sub-slice of internal buffer.
//
// Bytes are valid only until next call to any Decoder method.
func (d *Decoder) StrBytes() ([]byte, error) {
	v, err := d.str(value{raw: true})
	if err != nil {
		return nil, err
	}
	return v.buf, nil
}

// Str reads string.
func (d *Decoder) Str() (string, error) {
	s, err := d.StrBytes()
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (d *Decoder) strSlow(v value) (value, error) {
	for {
		c, err := d.byte()
		if err != nil {
			return value{}, errors.Wrap(err, "next")
		}
		switch c {
		case '"':
			// End of string.
			return v, nil
		case '\\':
			c, err := d.byte()
			if err != nil {
				return value{}, errors.Wrap(err, "next")
			}
			v, err = d.escapedChar(v, c)
			if err != nil {
				return v, errors.Wrap(err, "escape")
			}
		default:
			v = v.byte(c)
		}
	}
}

func (d *Decoder) escapedChar(v value, c byte) (value, error) {
	switch c {
	case 'u':
		r1, err := d.readU4()
		if err != nil {
			return value{}, errors.Wrap(err, "read u4")
		}
		if utf16.IsSurrogate(r1) {
			c, err := d.byte()
			if err != nil {
				return value{}, err
			}
			if c != '\\' {
				d.unread()
				return v.rune(r1), nil
			}
			c, err = d.byte()
			if err != nil {
				return value{}, err
			}
			if c != 'u' {
				return d.escapedChar(v.rune(r1), c)
			}
			r2, err := d.readU4()
			if err != nil {
				return value{}, err
			}
			combined := utf16.DecodeRune(r1, r2)
			if combined == '\uFFFD' {
				v = v.rune(r1).rune(r2)
			} else {
				v = v.rune(combined)
			}
		} else {
			v = v.rune(r1)
		}
	case '"':
		v = v.rune('"')
	case '\\':
		v = v.rune('\\')
	case '/':
		v = v.rune('/')
	case 'b':
		v = v.rune('\b')
	case 'f':
		v = v.rune('\f')
	case 'n':
		v = v.rune('\n')
	case 'r':
		v = v.rune('\r')
	case 't':
		v = v.rune('\t')
	default:
		return v, errors.Wrap(badToken(c), "bad escape: %w")
	}
	return v, nil
}

func (d *Decoder) readU4() (rune, error) {
	var v rune
	for i := 0; i < 4; i++ {
		c, err := d.byte()
		if err != nil {
			return 0, err
		}
		switch {
		case c >= '0' && c <= '9':
			v = v*16 + rune(c-'0')
		case c >= 'a' && c <= 'f':
			v = v*16 + rune(c-'a'+10)
		case c >= 'A' && c <= 'F':
			v = v*16 + rune(c-'A'+10)
		default:
			return 0, badToken(c)
		}
	}
	return v, nil
}

func appendRune(p []byte, r rune) []byte {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	return append(p, buf[:n]...)
}
