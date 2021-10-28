package jir

// ReadArray read array element, tells if the array has more element to read.
func (it *Iterator) ReadArray() (ret bool) {
	c := it.nextToken()
	switch c {
	case 'n':
		it.skipThreeBytes('u', 'l', 'l')
		return false // null
	case '[':
		c = it.nextToken()
		if c != ']' {
			it.unreadByte()
			return true
		}
		return false
	case ']':
		return false
	case ',':
		return true
	default:
		it.ReportError("ReadArray", "expect [ or , or ] or n, but found "+string([]byte{c}))
		return
	}
}

// ReadArrayCB read array with callback
func (it *Iterator) ReadArrayCB(callback func(*Iterator) bool) (ret bool) {
	c := it.nextToken()
	if c == '[' {
		if !it.incrementDepth() {
			return false
		}
		c = it.nextToken()
		if c != ']' {
			it.unreadByte()
			if !callback(it) {
				it.decrementDepth()
				return false
			}
			c = it.nextToken()
			for c == ',' {
				if !callback(it) {
					it.decrementDepth()
					return false
				}
				c = it.nextToken()
			}
			if c != ']' {
				it.ReportError("ReadArrayCB", "expect ] in the end, but found "+string([]byte{c}))
				it.decrementDepth()
				return false
			}
			return it.decrementDepth()
		}
		return it.decrementDepth()
	}
	if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	it.ReportError("ReadArrayCB", "expect [ or n, but found "+string([]byte{c}))
	return false
}
