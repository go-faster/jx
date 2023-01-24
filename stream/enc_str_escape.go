package stream

import (
	"io"
	"unicode/utf8"

	"github.com/go-faster/jx/internal/byteseq"
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
} // FIXME(tdakkota): re-use set

// StrEscape encodes string with html special characters escaping.
func (e *Encoder[W]) StrEscape(v string) bool {
	return e.comma() || strEscape(e, v)
}

// ByteStrEscape encodes string with html special characters escaping.
func (e *Encoder[W]) ByteStrEscape(v []byte) bool {
	return e.comma() || strEscape(e, v)
}

func strEscape[S byteseq.Byteseq, W io.Writer](e *Encoder[W], s S) bool {
	length := len(s)
	if e.w.writeByte('"') {
		return true
	}
	// Fast path, probably does not require escaping.
	i := 0
	for ; i < length; i++ {
		c := s[i]
		if c >= utf8.RuneSelf || !(htmlSafeSet[c]) {
			break
		}
	}
	if writeByteseq(&e.w, s[:i]) {
		return true
	}
	if i == length {
		return e.w.writeByte('"')
	}
	return strEscapeSlow[S](e, i, s, length)
}

func strEscapeSlow[S byteseq.Byteseq, W io.Writer](e *Encoder[W], i int, s S, valLen int) bool {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				if writeByteseq(&e.w, s[start:i]) {
					return true
				}
			}
			var ok bool
			switch b {
			case '\\', '"':
				ok = e.w.writeBytes('\\', b)
			case '\n':
				ok = e.w.writeBytes('\\', 'n')
			case '\r':
				ok = e.w.writeBytes('\\', 'r')
			case '\t':
				ok = e.w.writeBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				ok = e.w.writeString(`\u00`) || e.w.writeBytes(hexChars[b>>4], hexChars[b&0xF])
			}
			if !ok {
				return false
			}

			i++
			start = i
			continue
		}
		c, size := byteseq.DecodeRuneInByteseq(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				if writeByteseq(&e.w, s[start:i]) {
					return true
				}
			}
			if e.w.writeString(`\ufffd`) {
				return true
			}
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
				if writeByteseq(&e.w, s[start:i]) {
					return true
				}
			}
			if e.w.writeString(`\u202`) || e.w.writeByte(hexChars[c&0xF]) {
				return true
			}
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		if writeByteseq(&e.w, s[start:]) {
			return true
		}
	}
	return e.w.writeByte('"')
}
