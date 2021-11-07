package jx

import (
	"unicode/utf8"
)

// htmlSafeSet holds the value true if the ASCII character with the given
// array position can be safely represented inside a JSON string, embedded
// inside of HTML <script> tags, without any additional escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), the backslash character ("\"), HTML opening and closing
// tags ("<" and ">"), and the ampersand ("&").
var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

const hexChars = "0123456789abcdef"

// StrEscape encodes string with html special characters escaping.
func (e *Encoder) StrEscape(v string) {
	e.comma()
	length := len(v)
	e.buf = append(e.buf, '"')
	// Fast path, probably does not require escaping.
	i := 0
	for ; i < length; i++ {
		c := v[i]
		if c < utf8.RuneSelf && htmlSafeSet[c] {
			e.buf = append(e.buf, c)
		} else {
			break
		}
	}
	if i == length {
		e.buf = append(e.buf, '"')
		return
	}
	e.strEscape(i, v, length)
}

func (e *Encoder) strEscape(i int, v string, valLen int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := v[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				e.RawStr(v[start:i])
			}
			switch b {
			case '\\', '"':
				e.twoBytes('\\', b)
			case '\n':
				e.twoBytes('\\', 'n')
			case '\r':
				e.twoBytes('\\', 'r')
			case '\t':
				e.twoBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.RawStr(`\u00`)
				e.twoBytes(hexChars[b>>4], hexChars[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(v[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.RawStr(v[start:i])
			}
			e.RawStr(`\ufffd`)
			i++
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.RawStr(v[start:i])
			}
			e.RawStr(`\u202`)
			e.byte(hexChars[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(v) {
		e.RawStr(v[start:])
	}
	e.byte('"')
}

// Str encodes string without html escaping.
//
// Use StrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (e *Encoder) Str(v string) {
	e.comma()
	length := len(v)
	e.buf = append(e.buf, '"')
	// Fast path, without utf8 and escape support.
	i := 0
	for ; i < length; i++ {
		c := v[i]
		if c > 31 && c != '"' && c != '\\' && c < utf8.RuneSelf {
			e.buf = append(e.buf, c)
		} else {
			break
		}
	}
	if i == length {
		e.buf = append(e.buf, '"')
		return
	}
	e.strSlow(i, v, length)
}

func (e *Encoder) strSlow(i int, v string, length int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < length {
		if b := v[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				e.RawStr(v[start:i])
			}
			switch b {
			case '\\', '"':
				e.twoBytes('\\', b)
			case '\n':
				e.twoBytes('\\', 'n')
			case '\r':
				e.twoBytes('\\', 'r')
			case '\t':
				e.twoBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.RawStr(`\u00`)
				e.twoBytes(hexChars[b>>4], hexChars[b&0xF])
			}
			i++
			start = i
			continue
		}
		i++
		continue
	}
	if start < len(v) {
		e.RawStr(v[start:])
	}
	e.byte('"')
}
