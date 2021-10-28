package json

import (
	"fmt"
)

// ReadObject read one field from object.
// If object ended, returns empty string.
// Otherwise, returns the field name.
func (iter *Iterator) ReadObject() (ret string) {
	c := iter.nextToken()
	switch c {
	case 'n':
		iter.skipThreeBytes('u', 'l', 'l')
		return "" // null
	case '{':
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field := iter.ReadString()
			c = iter.nextToken()
			if c != ':' {
				iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			}
			return field
		}
		if c == '}' {
			return "" // end of object
		}
		iter.ReportError("ReadObject", `expect " after {, but found `+string([]byte{c}))
		return
	case ',':
		field := iter.ReadString()
		c = iter.nextToken()
		if c != ':' {
			iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
		}
		return field
	case '}':
		return "" // end of object
	default:
		iter.ReportError("ReadObject", fmt.Sprintf(`expect { or , or } or n, but found %s`, string([]byte{c})))
		return
	}
}

// ReadObjectCB read object with callback, the key is ascii only and field name not copied
func (iter *Iterator) ReadObjectCB(callback func(*Iterator, string) bool) bool {
	c := iter.nextToken()
	var field string
	if c == '{' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field = iter.ReadString()
			c = iter.nextToken()
			if c != ':' {
				iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			}
			if !callback(iter, field) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				field = iter.ReadString()
				c = iter.nextToken()
				if c != ':' {
					iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
				}
				if !callback(iter, field) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != '}' {
				iter.ReportError("ReadObjectCB", `object not ended with }`)
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		if c == '}' {
			return iter.decrementDepth()
		}
		iter.ReportError("ReadObjectCB", `expect " after {, but found `+string([]byte{c}))
		iter.decrementDepth()
		return false
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.ReportError("ReadObjectCB", `expect { or n, but found `+string([]byte{c}))
	return false
}

// ReadMapCB read map with callback, the key can be any string
func (iter *Iterator) ReadMapCB(callback func(*Iterator, string) bool) bool {
	c := iter.nextToken()
	if c == '{' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field := iter.ReadString()
			if iter.nextToken() != ':' {
				iter.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
				iter.decrementDepth()
				return false
			}
			if !callback(iter, field) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				field = iter.ReadString()
				if iter.nextToken() != ':' {
					iter.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
					iter.decrementDepth()
					return false
				}
				if !callback(iter, field) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != '}' {
				iter.ReportError("ReadMapCB", `object not ended with }`)
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		if c == '}' {
			return iter.decrementDepth()
		}
		iter.ReportError("ReadMapCB", `expect " after {, but found `+string([]byte{c}))
		iter.decrementDepth()
		return false
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
	return false
}
