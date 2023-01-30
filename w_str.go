package jx

import (
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
}

// Str encodes string without html escaping.
//
// Use StrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (w *Writer) Str(v string) bool {
	return writeStr(w, v)
}

// ByteStr encodes string without html escaping.
//
// Use ByteStrEscape to escape html, this is default for encoding/json and
// should be used by default for untrusted strings.
func (w *Writer) ByteStr(v []byte) bool {
	return writeStr(w, v)
}

func writeStr[S byteseq.Byteseq](w *Writer, v S) bool {
	if w.byte('"') {
		return true
	}

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
		return writeStreamByteseq(w, v) || w.byte('"')
	}
slow:
	w.Buf = append(w.Buf, v[:i]...)
	return strSlow[S](w, v[i:])
}

func strSlow[S byteseq.Byteseq](w *Writer, v S) bool {
	var i, start int
	// for the remaining parts, we process them char by char
	for i < len(v) {
		b := v[i]
		if safeSet[b] == 0 {
			i++
			continue
		}
		if start < i {
			if writeStreamByteseq(w, v[start:i]) {
				return true
			}
		}

		var fail bool
		switch b {
		case '\\', '"':
			fail = w.twoBytes('\\', b)
		case '\n':
			fail = w.twoBytes('\\', 'n')
		case '\r':
			fail = w.twoBytes('\\', 'r')
		case '\t':
			fail = w.twoBytes('\\', 't')
		default:
			// This encodes bytes < 0x20 except for \t, \n and \r.
			// If escapeHTML is set, it also escapes <, >, and &
			// because they can lead to security holes when
			// user-controlled strings are rendered into JSON
			// and served to some browsers.
			fail = w.rawStr(`\u00`) || w.twoBytes(hexChars[b>>4], hexChars[b&0xF])
		}
		if fail {
			return true
		}
		i++
		start = i
		continue
	}
	if start < len(v) {
		if writeStreamByteseq(w, v[start:]) {
			return true
		}
	}
	return w.byte('"')
}
