package jir

import "fmt"

// ReadNil reads a json object as nil and
// returns whether it's a nil or not
func (it *Iterator) ReadNil() (ret bool) {
	c := it.nextToken()
	if c == 'n' {
		it.skipThreeBytes('u', 'l', 'l') // null
		return true
	}
	it.unreadByte()
	return false
}

// ReadBool reads a json object as Bool
func (it *Iterator) ReadBool() (ret bool) {
	c := it.nextToken()
	if c == 't' {
		it.skipThreeBytes('r', 'u', 'e')
		return true
	}
	if c == 'f' {
		it.skipFourBytes('a', 'l', 's', 'e')
		return false
	}
	it.ReportError("ReadBool", "expect t or f, but found "+string([]byte{c}))
	return
}

// Skip skips a json object and positions to relatively the next json object.
func (it *Iterator) Skip() {
	c := it.nextToken()
	switch c {
	case '"':
		it.skipString()
	case 'n':
		it.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		it.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		it.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		it.unreadByte()
		it.ReadFloat32()
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		it.skipNumber()
	case '[':
		it.skipArray()
	case '{':
		it.skipObject()
	default:
		it.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
}

func (it *Iterator) skipFourBytes(b1, b2, b3, b4 byte) {
	if it.readByte() != b1 {
		it.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if it.readByte() != b2 {
		it.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if it.readByte() != b3 {
		it.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if it.readByte() != b4 {
		it.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
}

func (it *Iterator) skipThreeBytes(b1, b2, b3 byte) {
	if it.readByte() != b1 {
		it.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if it.readByte() != b2 {
		it.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if it.readByte() != b3 {
		it.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
}
