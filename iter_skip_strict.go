package jir

import (
	"fmt"
	"io"
)

func (it *Iterator) skipNumber() {
	if !it.trySkipNumber() {
		it.unreadByte()
		if it.Error != nil && it.Error != io.EOF {
			return
		}
		it.ReadFloat64()
		if it.Error != nil && it.Error != io.EOF {
			it.Error = nil
			it.ReadBigFloat()
		}
	}
}

func (it *Iterator) trySkipNumber() bool {
	dotFound := false
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case tDot:
			if dotFound {
				it.ReportError("validateNumber", `more than one dot found in number`)
				return true // already failed
			}
			if i+1 == it.tail {
				return false
			}
			c = it.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				it.ReportError("validateNumber", `missing digit after dot`)
				return true // already failed
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if it.head == i {
					return false // if - without following digits
				}
				it.head = i
				return true // must be valid
			}
			return false // may be invalid
		}
	}
	return false
}

func (it *Iterator) skipString() {
	if !it.trySkipString() {
		it.unreadByte()
		it.ReadString()
	}
}

func (it *Iterator) trySkipString() bool {
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		if c == '"' {
			it.head = i + 1
			return true // valid
		} else if c == '\\' {
			return false
		} else if c < ' ' {
			it.ReportError("trySkipString",
				fmt.Sprintf(`invalid control character found: %d`, c))
			return true // already failed
		}
	}
	return false
}

func (it *Iterator) skipObject() {
	it.unreadByte()
	it.ReadObjectCB(func(iter *Iterator, field string) bool {
		iter.Skip()
		return true
	})
}

func (it *Iterator) skipArray() {
	it.unreadByte()
	it.ReadArrayCB(func(iter *Iterator) bool {
		iter.Skip()
		return true
	})
}
