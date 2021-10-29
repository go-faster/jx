package jx

import (
	"fmt"
	"io"
	"unicode/utf16"

	"golang.org/x/xerrors"
)

// StrAppend reads string and appends it to byte slice.
func (it *Iter) StrAppend(b []byte) ([]byte, error) {
	v := value{
		buf: b,
		raw: false,
	}
	var err error
	if v, err = it.str(v); err != nil {
		return b, err
	}
	return v.buf, nil
}

type value struct {
	buf    []byte
	raw    bool // false forces buf reuse
	ignore bool
}

func (v value) direct(b []byte) value {
	if v.ignore {
		return v
	}
	if v.raw {
		return value{buf: b}
	}
	return v.append(b)
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

func (v value) append(b []byte) value {
	if v.ignore {
		return v
	}
	return value{
		buf: append(v.buf, b...),
		raw: v.raw,
	}
}

// UnexpectedTokenErr means that Token was unexpected while reading json.
type UnexpectedTokenErr struct {
	Token byte
}

func (e UnexpectedTokenErr) Error() string {
	return fmt.Sprintf("unexpected byte %d '%s'", e.Token, []byte{e.Token})
}

func badToken(c byte) error {
	return UnexpectedTokenErr{Token: c}
}

func (it *Iter) str(v value) (value, error) {
	if err := it.expectNext('"'); err != nil {
		return value{}, xerrors.Errorf("start: %w", err)
	}
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		if c == '\\' {
			// Character is escaped, fallback to slow path.
			break
		}
		if c == '"' {
			// End of string in fast path.
			str := it.buf[it.head:i]
			it.head = i + 1
			return v.direct(str), nil
		}
		if c < ' ' {
			return value{}, xerrors.Errorf("control character: %w", badToken(c))
		}
	}
	return it.strSlow(v)
}

// StrBytes returns string value as sub-slice of internal buffer.
//
// Buffer is valid only until next call to any Iter method.
func (it *Iter) StrBytes() ([]byte, error) {
	v, err := it.str(value{raw: true})
	if err != nil {
		return nil, err
	}
	return v.buf, nil
}

// Str reads string.
func (it *Iter) Str() (string, error) {
	s, err := it.StrBytes()
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (it *Iter) strSlow(v value) (value, error) {
	for {
		c, err := it.byte()
		if err == io.EOF {
			return value{}, io.ErrUnexpectedEOF
		}
		if err != nil {
			return value{}, xerrors.Errorf("next: %w", err)
		}
		switch c {
		case '"':
			// End of string.
			return v, nil
		case '\\':
			c, err := it.byte()
			if err == io.EOF {
				return value{}, io.ErrUnexpectedEOF
			}
			if err != nil {
				return value{}, xerrors.Errorf("next: %w", err)
			}
			v, err = it.escapedChar(v, c)
			if err != nil {
				return v, xerrors.Errorf("escape: %w", err)
			}
		default:
			v = v.byte(c)
		}
	}
}

func (it *Iter) escapedChar(v value, c byte) (value, error) {
	switch c {
	case 'u':
		r, err := it.readU4()
		if err != nil {
			return value{}, xerrors.Errorf("read u4: %w", err)
		}
		if utf16.IsSurrogate(r) {
			c, err := it.byte()
			if err == io.EOF {
				return value{}, io.ErrUnexpectedEOF
			}
			if err != nil {
				return value{}, err
			}
			if c != '\\' {
				it.unread()
				return v.rune(r), nil
			}
			c, err = it.byte()
			if err == io.EOF {
				return value{}, io.ErrUnexpectedEOF
			}
			if err != nil {
				return value{}, err
			}
			if c != 'u' {
				return it.escapedChar(v.rune(r), c)
			}
			r2, err := it.readU4()
			if err != nil {
				return value{}, err
			}
			combined := utf16.DecodeRune(r, r2)
			if combined == '\uFFFD' {
				v = v.rune(r).rune(r2)
			} else {
				v = v.rune(combined)
			}
		} else {
			v = v.rune(r)
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
		return v, xerrors.Errorf("bad escape: %w", badToken(c))
	}
	return v, nil
}

func (it *Iter) readU4() (rune, error) {
	var v rune
	for i := 0; i < 4; i++ {
		c, err := it.byte()
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		}
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

//nolint:unused,deadcode,varcheck
const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	maxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
	runeError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
)

func appendRune(p []byte, r rune) []byte {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		return append(p, byte(r))
	case i <= rune2Max:
		return append(p,
			t2|byte(r>>6),
			tx|byte(r)&maskx,
		)
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = runeError
		fallthrough
	case i <= rune3Max:
		return append(p,
			t3|byte(r>>12),
			tx|byte(r>>6)&maskx,
			tx|byte(r)&maskx,
		)
	default:
		return append(p,
			t4|byte(r>>18),
			tx|byte(r>>12)&maskx,
			tx|byte(r>>6)&maskx,
			tx|byte(r)&maskx,
		)
	}
}
