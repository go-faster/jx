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
			field := it.String()
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
		field := it.String()
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

// Object read object, calling f on each field.
func (it *Iterator) Object(f func(i *Iterator, key string) bool) bool {
	c := it.nextToken()
	if c == '{' {
		if !it.incrementDepth() {
			return false
		}
		c = it.nextToken()
		if c == '"' {
			it.unreadByte()
			key := it.String()
			c = it.nextToken()
			if c != ':' {
				it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
			}
			if !f(it, key) {
				it.decrementDepth()
				return false
			}
			c = it.nextToken()
			for c == ',' {
				key = it.String()
				c = it.nextToken()
				if c != ':' {
					it.ReportError("Field", "expect : after object field, but found "+string([]byte{c}))
				}
				if !f(it, key) {
					it.decrementDepth()
					return false
				}
				c = it.nextToken()
			}
			if c != '}' {
				it.ReportError("Object", `object not ended with }`)
				it.decrementDepth()
				return false
			}
			return it.decrementDepth()
		}
		if c == '}' {
			return it.decrementDepth()
		}
		it.ReportError("Object", `expect " after {, but found `+string([]byte{c}))
		it.decrementDepth()
		return false
	}
	if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	it.ReportError("Object", `expect { or n, but found `+string([]byte{c}))
	return false
}
