package stream

import (
	"io"

	"github.com/go-faster/jx/internal/byteseq"
)

const hexChars = "0123456789abcdef"

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
} // FIXME(tdakkota): re-use set

// Str encodes string without html escaping.
//
// Use StrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (e *Encoder[W]) Str(v string) bool {
	return e.comma() || writeStr(e, v)
}

// ByteStr encodes string without html escaping.
//
// Use ByteStrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (e *Encoder[W]) ByteStr(v []byte) bool {
	return e.comma() || writeStr(e, v)
}

func writeStr[S byteseq.Byteseq, W io.Writer](e *Encoder[W], s S) bool {
	if e.w.writeByte('"') {
		return true
	}

	// Fast path, without utf8 and escape support.
	var (
		i      = 0
		length = len(s)
		c      byte
	)
	for i, c = range []byte(s) {
		if safeSet[c] != 0 {
			goto slow
		}
	}
	if i == length-1 {
		return writeByteseq(&e.w, s) || e.w.writeByte('"')
	}
slow:
	return writeByteseq(&e.w, s[:i]) || strSlow[S](e, s[i:])
}

func strSlow[S byteseq.Byteseq, W io.Writer](e *Encoder[W], s S) bool {
	var i, start int
	// for the remaining parts, we process them char by char
	for i < len(s) {
		b := s[i]
		if safeSet[b] == 0 {
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
			ok = e.w.writeString(`\u00`) ||
				e.w.writeBytes(hexChars[b>>4], hexChars[b&0xF])
		}
		if !ok {
			return false
		}

		i++
		start = i
		continue
	}
	if start < len(s) {
		if writeByteseq(&e.w, s[start:]) {
			return false
		}
	}
	return e.w.writeByte('"')
}
