package jir

import (
	"fmt"
	"unicode/utf16"
)

// String reads string.
func (it *Iterator) String() (s string) {
	c := it.nextToken()
	if c == '"' {
		for i := it.head; i < it.tail; i++ {
			c := it.buf[i]
			if c == '"' {
				s = string(it.buf[it.head:i])
				it.head = i + 1
				return s
			} else if c == '\\' {
				break
			} else if c < ' ' {
				it.ReportError("String",
					fmt.Sprintf(`invalid control character found: %d`, c))
				return
			}
		}
		return it.readStringSlowPath()
	} else if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l')
		return ""
	}
	it.ReportError("String", `expects " or n, but found `+string([]byte{c}))
	return ""
}

func (it *Iterator) readStringSlowPath() (ret string) {
	var str []byte
	var c byte
	for it.Error == nil {
		c = it.readByte()
		if c == '"' {
			return string(str)
		}
		if c == '\\' {
			c = it.readByte()
			str = it.readEscapedChar(c, str)
		} else {
			str = append(str, c)
		}
	}
	it.ReportError("readStringSlowPath", "unexpected end of input")
	return
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

// StringAsSlice read string from iterator without copying into string form.
// The []byte can not be kept, as it will change after next iterator call.
func (it *Iterator) StringAsSlice() (ret []byte) {
	c := it.nextToken()
	if c == '"' {
		for i := it.head; i < it.tail; i++ {
			// require ascii string and no escape
			// for: field name, base64, number
			if it.buf[i] == '"' {
				// fast path: reuse the underlying buffer
				ret = it.buf[it.head:i]
				it.head = i + 1
				return ret
			}
		}
		readLen := it.tail - it.head
		copied := make([]byte, readLen, readLen*2)
		copy(copied, it.buf[it.head:it.tail])
		it.head = it.tail
		for it.Error == nil {
			c := it.readByte()
			if c == '"' {
				return copied
			}
			copied = append(copied, c)
		}
		return copied
	}
	it.ReportError("StringAsSlice", `expects " or n, but found `+string([]byte{c}))
	return
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
