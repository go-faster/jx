package jir

// Elem reads array element and reports whether array has more
// elements to read.
func (it *Iterator) Elem() (ret bool) {
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
		it.ReportError("Elem", "expect [ or , or ] or n, but found "+string([]byte{c}))
		return
	}
}

// Array reads array and call f on each element.
func (it *Iterator) Array(f func(i *Iterator) bool) (ret bool) {
	c := it.nextToken()
	if c == '[' {
		if !it.incrementDepth() {
			return false
		}
		c = it.nextToken()
		if c != ']' {
			it.unreadByte()
			if !f(it) {
				it.decrementDepth()
				return false
			}
			c = it.nextToken()
			for c == ',' {
				if !f(it) {
					it.decrementDepth()
					return false
				}
				c = it.nextToken()
			}
			if c != ']' {
				it.ReportError("Array", "expect ] in the end, but found "+string([]byte{c}))
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
	it.ReportError("Array", "expect [ or n, but found "+string([]byte{c}))
	return false
}
