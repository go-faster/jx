package jir

import (
	"fmt"
	"unicode/utf16"
)

// StrAppend reads string and appends it to byte slice.
func (it *Iterator) StrAppend(b []byte) []byte {
	return it.strBytes(b)
}

func (it *Iterator) strBytes(b []byte) []byte {
	c := it.nextToken()
	if c == '"' {
		for i := it.head; i < it.tail; i++ {
			c := it.buf[i]
			if c == '\\' {
				break // escaped, fallback to slow path
			}
			if c == '"' {
				// End of string in fast path.
				str := it.buf[it.head:i]
				it.head = i + 1
				if b == nil {
					// Returning str directly if no b.
					return str
				}
				return append(b, str...)
			}
			if c < ' ' {
				it.ReportError("Str", fmt.Sprintf(`invalid control character found: %d`, c))
				return b
			}
		}
		return it.readStringSlowPath(b)
	}
	it.ReportError("Str", `expects " or n, but found `+string([]byte{c}))
	return nil
}

// StrBytes returns string value as sub-slice of internal buffer.
//
// Buffer is valid only until next call to any Iterator method.
func (it *Iterator) StrBytes() []byte {
	return it.strBytes(nil)
}

// Str reads string.
func (it *Iterator) Str() string {
	if s := it.StrBytes(); s != nil {
		return string(s)
	}
	return ""
}

func (it *Iterator) readStringSlowPath(b []byte) []byte {
	for it.Error == nil {
		c := it.readByte()
		if c == '"' {
			return b // end of string
		}
		if c == '\\' {
			if str := it.readEscapedChar(it.readByte(), b); str == nil {
				return nil
			} else {
				b = str
			}
		} else {
			b = append(b, c)
		}
	}
	it.ReportError("readStringSlowPath", "unexpected end of input")
	return nil
}

func (it *Iterator) readEscapedChar(c byte, str []byte) []byte {
	switch c {
	case 'u':
		r := it.readU4()
		if utf16.IsSurrogate(r) {
			c = it.readByte()
			if it.Error != nil {
				return nil
			}
			if c != '\\' {
				it.unreadByte()
				str = appendRune(str, r)
				return str
			}
			c = it.readByte()
			if it.Error != nil {
				return nil
			}
			if c != 'u' {
				str = appendRune(str, r)
				return it.readEscapedChar(c, str)
			}
			r2 := it.readU4()
			if it.Error != nil {
				return nil
			}
			combined := utf16.DecodeRune(r, r2)
			if combined == '\uFFFD' {
				str = appendRune(str, r)
				str = appendRune(str, r2)
			} else {
				str = appendRune(str, combined)
			}
		} else {
			str = appendRune(str, r)
		}
	case '"':
		str = append(str, '"')
	case '\\':
		str = append(str, '\\')
	case '/':
		str = append(str, '/')
	case 'b':
		str = append(str, '\b')
	case 'f':
		str = append(str, '\f')
	case 'n':
		str = append(str, '\n')
	case 'r':
		str = append(str, '\r')
	case 't':
		str = append(str, '\t')
	default:
		it.ReportError("readEscapedChar",
			`invalid escape char after \`)
		return nil
	}
	return str
}

func (it *Iterator) readU4() (ret rune) {
	for i := 0; i < 4; i++ {
		c := it.readByte()
		if it.Error != nil {
			return
		}
		if c >= '0' && c <= '9' {
			ret = ret*16 + rune(c-'0')
		} else if c >= 'a' && c <= 'f' {
			ret = ret*16 + rune(c-'a'+10)
		} else if c >= 'A' && c <= 'F' {
			ret = ret*16 + rune(c-'A'+10)
		} else {
			it.ReportError("readU4", "expects 0~9 or a~f, but found "+string([]byte{c}))
			return
		}
	}
	return ret
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
		p = append(p, byte(r))
		return p
	case i <= rune2Max:
		p = append(p, t2|byte(r>>6))
		p = append(p, tx|byte(r)&maskx)
		return p
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = runeError
		fallthrough
	case i <= rune3Max:
		p = append(p, t3|byte(r>>12))
		p = append(p, tx|byte(r>>6)&maskx)
		p = append(p, tx|byte(r)&maskx)
		return p
	default:
		p = append(p, t4|byte(r>>18))
		p = append(p, tx|byte(r>>12)&maskx)
		p = append(p, tx|byte(r>>6)&maskx)
		p = append(p, tx|byte(r)&maskx)
		return p
	}
}
