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

const hexChars = "0123456789abcdef"

// StrEscape encodes string with html special characters escaping.
func (w *Writer) StrEscape(v string) {
	strEscape(w, v)
}

// ByteStrEscape encodes string with html special characters escaping.
func (w *Writer) ByteStrEscape(v []byte) {
	strEscape(w, v)
}

func strEscape[T byteseq](w *Writer, v T) {
	length := len(v)
	w.Buf = append(w.Buf, '"')
	// Fast path, probably does not require escaping.
	i := 0
	for ; i < length; i++ {
		c := v[i]
		if c >= utf8.RuneSelf || !(htmlSafeSet[c]) {
			break
		}
	}
	w.Buf = append(w.Buf, v[:i]...)
	if i == length {
		w.Buf = append(w.Buf, '"')
		return
	}
	strEscapeSlow[T](w, i, v, length)
}

func strEscapeSlow[T byteseq](w *Writer, i int, v T, valLen int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := v[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				w.Buf = append(w.Buf, v[start:i]...)
			}
			switch b {
			case '\\', '"':
				w.twoBytes('\\', b)
			case '\n':
				w.twoBytes('\\', 'n')
			case '\r':
				w.twoBytes('\\', 'r')
			case '\t':
				w.twoBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.rawStr(`\u00`)
				w.twoBytes(hexChars[b>>4], hexChars[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := decodeRuneInByteseq(v[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.Buf = append(w.Buf, v[start:i]...)
			}
			w.rawStr(`\ufffd`)
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
				w.Buf = append(w.Buf, v[start:i]...)
			}
			w.rawStr(`\u202`)
			w.byte(hexChars[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(v) {
		w.Buf = append(w.Buf, v[start:]...)
	}
	w.byte('"')
}

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [256]byte{
	// First 31 characters.
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	'"':  1,
	'\\': 1,
}

// Str encodes string without html escaping.
//
// Use StrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (w *Writer) Str(v string) {
	writeStr(w, v)
}

// ByteStr encodes string without html escaping.
//
// Use ByteStrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (w *Writer) ByteStr(v []byte) {
	writeStr(w, v)
}

func writeStr[T byteseq](w *Writer, v T) {
	w.Buf = append(w.Buf, '"')

	// Fast path, without utf8 and escape support.
	var (
		i      = 0
		length = len(v)
		c      byte
	)
	for i, c = range []byte(v) {
		if safeSet[c] != 0 {
			goto slow
		}
	}
	if i == length-1 {
		w.Buf = append(w.Buf, v...)
		w.Buf = append(w.Buf, '"')
		return
	}
slow:
	w.Buf = append(w.Buf, v[:i]...)
	strSlow[T](w, v[i:])
}

func strSlow[T byteseq](w *Writer, v T) {
	var i, start int
	// for the remaining parts, we process them char by char
	for i < len(v) {
		b := v[i]
		if safeSet[b] == 0 {
			i++
			continue
		}
		if start < i {
			w.Buf = append(w.Buf, v[start:i]...)
		}
		switch b {
		case '\\', '"':
			w.twoBytes('\\', b)
		case '\n':
			w.twoBytes('\\', 'n')
		case '\r':
			w.twoBytes('\\', 'r')
		case '\t':
			w.twoBytes('\\', 't')
		default:
			// This encodes bytes < 0x20 except for \t, \n and \r.
			// If escapeHTML is set, it also escapes <, >, and &
			// because they can lead to security holes when
			// user-controlled strings are rendered into JSON
			// and served to some browsers.
			w.rawStr(`\u00`)
			w.twoBytes(hexChars[b>>4], hexChars[b&0xF])
		}
		i++
		start = i
		continue
	}
	if start < len(v) {
		w.Buf = append(w.Buf, v[start:]...)
	}
	w.byte('"')
}
