package jir

import (
	"math"
	"strconv"
)

var intDigits []int8

const uint32SafeToMultiply10 = uint32(0xffffffff)/10 - 1
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff)/10 - 1
const maxFloat64 = 1<<53 - 1

func init() {
	intDigits = make([]int8, 256)
	for i := 0; i < len(intDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

// Uint read uint.
func (it *Iterator) Uint() uint {
	if strconv.IntSize == 32 {
		return uint(it.Uint32())
	}
	return uint(it.Uint64())
}

// Int reads integer.
func (it *Iterator) Int() int {
	if strconv.IntSize == 32 {
		return int(it.Int32())
	}
	return int(it.Int64())
}

// Int8 read int8
func (it *Iterator) Int8() (ret int8) {
	c := it.nextToken()
	if c == '-' {
		val := it.readUint32(it.readByte())
		if val > math.MaxInt8+1 {
			it.ReportError("Int8", "overflow: "+strconv.FormatInt(int64(val), 10))
			return
		}
		return -int8(val)
	}
	val := it.readUint32(c)
	if val > math.MaxInt8 {
		it.ReportError("Int8", "overflow: "+strconv.FormatInt(int64(val), 10))
		return
	}
	return int8(val)
}

// Uint8 read uint8
func (it *Iterator) Uint8() (ret uint8) {
	val := it.readUint32(it.nextToken())
	if val > math.MaxUint8 {
		it.ReportError("Uint8", "overflow: "+strconv.FormatInt(int64(val), 10))
		return
	}
	return uint8(val)
}

// Int16 read int16
func (it *Iterator) Int16() (ret int16) {
	c := it.nextToken()
	if c == '-' {
		val := it.readUint32(it.readByte())
		if val > math.MaxInt16+1 {
			it.ReportError("Int16", "overflow: "+strconv.FormatInt(int64(val), 10))
			return
		}
		return -int16(val)
	}
	val := it.readUint32(c)
	if val > math.MaxInt16 {
		it.ReportError("Int16", "overflow: "+strconv.FormatInt(int64(val), 10))
		return
	}
	return int16(val)
}

// Uint16 read uint16
func (it *Iterator) Uint16() (ret uint16) {
	val := it.readUint32(it.nextToken())
	if val > math.MaxUint16 {
		it.ReportError("Uint16", "overflow: "+strconv.FormatInt(int64(val), 10))
		return
	}
	return uint16(val)
}

// Int32 read int32
func (it *Iterator) Int32() (ret int32) {
	c := it.nextToken()
	if c == '-' {
		val := it.readUint32(it.readByte())
		if val > math.MaxInt32+1 {
			it.ReportError("Int32", "overflow: "+strconv.FormatInt(int64(val), 10))
			return
		}
		return -int32(val)
	}
	val := it.readUint32(c)
	if val > math.MaxInt32 {
		it.ReportError("Int32", "overflow: "+strconv.FormatInt(int64(val), 10))
		return
	}
	return int32(val)
}

// Uint32 read uint32
func (it *Iterator) Uint32() (ret uint32) {
	return it.readUint32(it.nextToken())
}

func (it *Iterator) readUint32(c byte) (ret uint32) {
	ind := intDigits[c]
	if ind == 0 {
		it.assertInteger()
		return 0 // single zero
	}
	if ind == invalidCharForNumber {
		it.ReportError("readUint32", "unexpected character: "+string([]byte{byte(ind)}))
		return
	}
	value := uint32(ind)
	if it.tail-it.head > 10 {
		i := it.head
		ind2 := intDigits[it.buf[i]]
		if ind2 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value
		}
		i++
		ind3 := intDigits[it.buf[i]]
		if ind3 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*10 + uint32(ind2)
		}
		//iter.head = i + 1
		//value = value * 100 + uint32(ind2) * 10 + uint32(ind3)
		i++
		ind4 := intDigits[it.buf[i]]
		if ind4 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*100 + uint32(ind2)*10 + uint32(ind3)
		}
		i++
		ind5 := intDigits[it.buf[i]]
		if ind5 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*1000 + uint32(ind2)*100 + uint32(ind3)*10 + uint32(ind4)
		}
		i++
		ind6 := intDigits[it.buf[i]]
		if ind6 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*10000 + uint32(ind2)*1000 + uint32(ind3)*100 + uint32(ind4)*10 + uint32(ind5)
		}
		i++
		ind7 := intDigits[it.buf[i]]
		if ind7 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*100000 + uint32(ind2)*10000 + uint32(ind3)*1000 + uint32(ind4)*100 + uint32(ind5)*10 + uint32(ind6)
		}
		i++
		ind8 := intDigits[it.buf[i]]
		if ind8 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*1000000 + uint32(ind2)*100000 + uint32(ind3)*10000 + uint32(ind4)*1000 + uint32(ind5)*100 + uint32(ind6)*10 + uint32(ind7)
		}
		i++
		ind9 := intDigits[it.buf[i]]
		value = value*10000000 + uint32(ind2)*1000000 + uint32(ind3)*100000 + uint32(ind4)*10000 + uint32(ind5)*1000 + uint32(ind6)*100 + uint32(ind7)*10 + uint32(ind8)
		it.head = i
		if ind9 == invalidCharForNumber {
			it.assertInteger()
			return value
		}
	}
	for {
		for i := it.head; i < it.tail; i++ {
			ind = intDigits[it.buf[i]]
			if ind == invalidCharForNumber {
				it.head = i
				it.assertInteger()
				return value
			}
			if value > uint32SafeToMultiply10 {
				value2 := (value << 3) + (value << 1) + uint32(ind)
				if value2 < value {
					it.ReportError("readUint32", "overflow")
					return
				}
				value = value2
				continue
			}
			value = (value << 3) + (value << 1) + uint32(ind)
		}
		if !it.loadMore() {
			it.assertInteger()
			return value
		}
	}
}

// Int64 read int64
func (it *Iterator) Int64() (ret int64) {
	c := it.nextToken()
	if c == '-' {
		val := it.readUint64(it.readByte())
		if val > math.MaxInt64+1 {
			it.ReportError("Int64", "overflow: "+strconv.FormatUint(uint64(val), 10))
			return
		}
		return -int64(val)
	}
	val := it.readUint64(c)
	if val > math.MaxInt64 {
		it.ReportError("Int64", "overflow: "+strconv.FormatUint(uint64(val), 10))
		return
	}
	return int64(val)
}

// Uint64 read uint64
func (it *Iterator) Uint64() uint64 {
	return it.readUint64(it.nextToken())
}

func (it *Iterator) readUint64(c byte) (ret uint64) {
	ind := intDigits[c]
	if ind == 0 {
		it.assertInteger()
		return 0 // single zero
	}
	if ind == invalidCharForNumber {
		it.ReportError("readUint64", "unexpected character: "+string([]byte{byte(ind)}))
		return
	}
	value := uint64(ind)
	if it.tail-it.head > 10 {
		i := it.head
		ind2 := intDigits[it.buf[i]]
		if ind2 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value
		}
		i++
		ind3 := intDigits[it.buf[i]]
		if ind3 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*10 + uint64(ind2)
		}
		//iter.head = i + 1
		//value = value * 100 + uint32(ind2) * 10 + uint32(ind3)
		i++
		ind4 := intDigits[it.buf[i]]
		if ind4 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*100 + uint64(ind2)*10 + uint64(ind3)
		}
		i++
		ind5 := intDigits[it.buf[i]]
		if ind5 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*1000 + uint64(ind2)*100 + uint64(ind3)*10 + uint64(ind4)
		}
		i++
		ind6 := intDigits[it.buf[i]]
		if ind6 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*10000 + uint64(ind2)*1000 + uint64(ind3)*100 + uint64(ind4)*10 + uint64(ind5)
		}
		i++
		ind7 := intDigits[it.buf[i]]
		if ind7 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*100000 + uint64(ind2)*10000 + uint64(ind3)*1000 + uint64(ind4)*100 + uint64(ind5)*10 + uint64(ind6)
		}
		i++
		ind8 := intDigits[it.buf[i]]
		if ind8 == invalidCharForNumber {
			it.head = i
			it.assertInteger()
			return value*1000000 + uint64(ind2)*100000 + uint64(ind3)*10000 + uint64(ind4)*1000 + uint64(ind5)*100 + uint64(ind6)*10 + uint64(ind7)
		}
		i++
		ind9 := intDigits[it.buf[i]]
		value = value*10000000 + uint64(ind2)*1000000 + uint64(ind3)*100000 + uint64(ind4)*10000 + uint64(ind5)*1000 + uint64(ind6)*100 + uint64(ind7)*10 + uint64(ind8)
		it.head = i
		if ind9 == invalidCharForNumber {
			it.assertInteger()
			return value
		}
	}
	for {
		for i := it.head; i < it.tail; i++ {
			ind = intDigits[it.buf[i]]
			if ind == invalidCharForNumber {
				it.head = i
				it.assertInteger()
				return value
			}
			if value > uint64SafeToMultiple10 {
				value2 := (value << 3) + (value << 1) + uint64(ind)
				if value2 < value {
					it.ReportError("readUint64", "overflow")
					return
				}
				value = value2
				continue
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
		if !it.loadMore() {
			it.assertInteger()
			return value
		}
	}
}

func (it *Iterator) assertInteger() {
	if it.head < it.tail && it.buf[it.head] == '.' {
		it.ReportError("assertInteger", "can not decode float as int")
	}
}
