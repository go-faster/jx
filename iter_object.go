package jir

import (
	"fmt"
)

// ReadObject read one field from object.
// If object ended, returns empty string.
// Otherwise, returns the field name.
func (it *Iterator) ReadObject() (ret string) {
	c := it.nextToken()
	switch c {
	case 'n':
		it.skipThreeBytes('u', 'l', 'l')
		return "" // null
	case '{':
		c = it.nextToken()
		if c == '"' {
			it.unreadByte()
			field := it.ReadString()
			c = it.nextToken()
			if c != ':' {
				it.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			}
			return field
		}
		if c == '}' {
			return "" // end of object
		}
		it.ReportError("ReadObject", `expect " after {, but found `+string([]byte{c}))
		return
	case ',':
		field := it.ReadString()
		c = it.nextToken()
		if c != ':' {
			it.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
		}
		return field
	case '}':
		return "" // end of object
	default:
		it.ReportError("ReadObject", fmt.Sprintf(`expect { or , or } or n, but found %s`, string([]byte{c})))
		return
	}
}

// ReadObjectCB read object with callback, the key is ascii only and field name not copied
func (it *Iterator) ReadObjectCB(callback func(*Iterator, string) bool) bool {
	c := it.nextToken()
	var field string
	if c == '{' {
		if !it.incrementDepth() {
			return false
		}
		c = it.nextToken()
		if c == '"' {
			it.unreadByte()
			field = it.ReadString()
			c = it.nextToken()
			if c != ':' {
				it.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			}
			if !callback(it, field) {
				it.decrementDepth()
				return false
			}
			c = it.nextToken()
			for c == ',' {
				field = it.ReadString()
				c = it.nextToken()
				if c != ':' {
					it.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
				}
				if !callback(it, field) {
					it.decrementDepth()
					return false
				}
				c = it.nextToken()
			}
			if c != '}' {
				it.ReportError("ReadObjectCB", `object not ended with }`)
				it.decrementDepth()
				return false
			}
			return it.decrementDepth()
		}
		if c == '}' {
			return it.decrementDepth()
		}
		it.ReportError("ReadObjectCB", `expect " after {, but found `+string([]byte{c}))
		it.decrementDepth()
		return false
	}
	if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	it.ReportError("ReadObjectCB", `expect { or n, but found `+string([]byte{c}))
	return false
}

// ReadMapCB read map with callback, the key can be any string
func (it *Iterator) ReadMapCB(callback func(*Iterator, string) bool) bool {
	c := it.nextToken()
	if c == '{' {
		if !it.incrementDepth() {
			return false
		}
		c = it.nextToken()
		if c == '"' {
			it.unreadByte()
			field := it.ReadString()
			if it.nextToken() != ':' {
				it.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
				it.decrementDepth()
				return false
			}
			if !callback(it, field) {
				it.decrementDepth()
				return false
			}
			c = it.nextToken()
			for c == ',' {
				field = it.ReadString()
				if it.nextToken() != ':' {
					it.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
					it.decrementDepth()
					return false
				}
				if !callback(it, field) {
					it.decrementDepth()
					return false
				}
				c = it.nextToken()
			}
			if c != '}' {
				it.ReportError("ReadMapCB", `object not ended with }`)
				it.decrementDepth()
				return false
			}
			return it.decrementDepth()
		}
		if c == '}' {
			return it.decrementDepth()
		}
		it.ReportError("ReadMapCB", `expect " after {, but found `+string([]byte{c}))
		it.decrementDepth()
		return false
	}
	if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	it.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
	return false
}
