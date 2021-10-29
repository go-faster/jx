package jir

import (
	"fmt"
)

// Field read one field from object.
// If object ended, returns empty string.
// Otherwise, returns the field name.
func (it *Iterator) Field() (ret string) {
	c := it.nextToken()
	switch c {
	case 'n':
		it.skipThreeBytes('u', 'l', 'l')
		return "" // null
	case '{':
		c = it.nextToken()
		if c == '"' {
			it.unreadByte()
			field := it.Str()
			c = it.nextToken()
			if c != ':' {
				it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
			}
			return field
		}
		if c == '}' {
			return "" // end of object
		}
		it.ReportError("Field", `expect " after {, but found `+string([]byte{c}))
		return
	case ',':
		field := it.Str()
		c = it.nextToken()
		if c != ':' {
			it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
		}
		return field
	case '}':
		return "" // end of object
	default:
		it.ReportError("Field", fmt.Sprintf(`expect { or , or } or n, but found %s`, string([]byte{c})))
		return
	}
}

func (it *Iterator) object(f func(i *Iterator, key []byte) bool) bool {
	if it.buf == nil {
		it.buf = make([]byte, 0, 64)
	}
	// Use it.buf to hold keys.
	// Rest back on exit.
	n := len(it.buf)
	defer func() { it.buf = it.buf[:n] }()

	c := it.nextToken()
	if c != '{' {
		it.ReportError("Object", `expect { or n, but found `+string([]byte{c}))
		return false
	}
	if !it.incrementDepth() {
		return false
	}

	c = it.nextToken()
	if c == '}' {
		return it.decrementDepth()
	}
	it.unreadByte()

	j := len(it.buf)
	if str := it.strBytes(it.buf); str == nil {
		return false
	} else {
		it.buf = str
	}
	k := it.buf[j:]

	c = it.nextToken()
	if c != ':' {
		it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
		return false
	}
	if !f(it, k) {
		it.decrementDepth()
		return false
	}

	// Drop k.
	it.buf = it.buf[:j]

	c = it.nextToken()
	for c == ',' {
		// Expand buf for k.
		j := len(it.buf)
		if str := it.strBytes(it.buf); str == nil {
			return false
		} else {
			it.buf = str
		}
		k := it.buf[j:]

		c = it.nextToken()
		if c != ':' {
			it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
		}
		if !f(it, k) {
			it.decrementDepth()
			return false
		}

		// Drop k.
		it.buf = it.buf[:j]

		c = it.nextToken()
	}
	if c != '}' {
		it.ReportError("Object", `object not ended with }`)
		it.decrementDepth()
		return false
	}
	return it.decrementDepth()
}

// Object read object, calling f on each field.
func (it *Iterator) Object(f func(i *Iterator, key string) bool) bool {
	return it.object(func(i *Iterator, key []byte) bool {
		return f(i, string(key))
	})
}
