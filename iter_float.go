package jir

import (
	"encoding/json"
	"io"
	"math/big"
	"strconv"
	"strings"
	"unsafe"
)

var floatDigits []int8

const invalidCharForNumber = int8(-1)
const endOfNumber = int8(-2)
const dotInNumber = int8(-3)

func init() {
	floatDigits = make([]int8, 256)
	for i := 0; i < len(floatDigits); i++ {
		floatDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		floatDigits[i] = i - int8('0')
	}
	floatDigits[','] = endOfNumber
	floatDigits[']'] = endOfNumber
	floatDigits['}'] = endOfNumber
	floatDigits[' '] = endOfNumber
	floatDigits['\t'] = endOfNumber
	floatDigits['\n'] = endOfNumber
	floatDigits[tDot] = dotInNumber
}

// ReadBigFloat read big.Float
func (it *Iterator) ReadBigFloat() (ret *big.Float) {
	str := it.readNumberAsString()
	if it.Error != nil && it.Error != io.EOF {
		return nil
	}
	prec := 64
	if len(str) > prec {
		prec = len(str)
	}
	val, _, err := big.ParseFloat(str, 10, uint(prec), big.ToZero)
	if err != nil {
		it.Error = err
		return nil
	}
	return val
}

// ReadBigInt read big.Int
func (it *Iterator) ReadBigInt() (ret *big.Int) {
	str := it.readNumberAsString()
	if it.Error != nil && it.Error != io.EOF {
		return nil
	}
	ret = big.NewInt(0)
	var success bool
	ret, success = ret.SetString(str, 10)
	if !success {
		it.ReportError("ReadBigInt", "invalid big int")
		return nil
	}
	return ret
}

//ReadFloat32 read float32
func (it *Iterator) ReadFloat32() (ret float32) {
	c := it.nextToken()
	if c == '-' {
		return -it.readPositiveFloat32()
	}
	it.unreadByte()
	return it.readPositiveFloat32()
}

func (it *Iterator) readPositiveFloat32() (ret float32) {
	i := it.head
	// First char.
	if i == it.tail {
		return it.readFloat32SlowPath()
	}
	c := it.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return it.readFloat32SlowPath()
	case endOfNumber:
		it.ReportError("readFloat32", "empty number")
		return
	case dotInNumber:
		it.ReportError("readFloat32", "leading dot is invalid")
		return
	case 0:
		if i == it.tail {
			return it.readFloat32SlowPath()
		}
		c = it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			it.ReportError("readFloat32", "leading zero is invalid")
			return
		}
	}
	value := uint64(ind)
	// Chars before dot.
NonDecimalLoop:
	for ; i < it.tail; i++ {
		c = it.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return it.readFloat32SlowPath()
		case endOfNumber:
			it.head = i
			return float32(value)
		case dotInNumber:
			break NonDecimalLoop
		}
		if value > uint64SafeToMultiple10 {
			return it.readFloat32SlowPath()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == tDot {
		i++
		decimalPlaces := 0
		if i == it.tail {
			return it.readFloat32SlowPath()
		}
		for ; i < it.tail; i++ {
			c = it.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					it.head = i
					return float32(float64(value) / float64(pow10[decimalPlaces]))
				}
				// too many decimal places
				return it.readFloat32SlowPath()
			case invalidCharForNumber, dotInNumber:
				return it.readFloat32SlowPath()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return it.readFloat32SlowPath()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return it.readFloat32SlowPath()
}

func (it *Iterator) readNumberAsString() (ret string) {
	strBuf := [16]byte{}
	str := strBuf[:]
Load:
	for {
		for i := it.head; i < it.tail; i++ {
			c := it.buf[i]
			switch c {
			case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			default:
				it.head = i
				break Load
			}
		}
		if !it.loadMore() {
			break
		}
	}
	if it.Error != nil && it.Error != io.EOF {
		return
	}
	if len(str) == 0 {
		it.ReportError("readNumberAsString", "invalid number")
	}
	return *(*string)(unsafe.Pointer(&str))
}

func (it *Iterator) readFloat32SlowPath() (ret float32) {
	str := it.readNumberAsString()
	if it.Error != nil && it.Error != io.EOF {
		return
	}
	errMsg := validateFloat(str)
	if errMsg != "" {
		it.ReportError("readFloat32SlowPath", errMsg)
		return
	}
	val, err := strconv.ParseFloat(str, 32)
	if err != nil {
		it.Error = err
		return
	}
	return float32(val)
}

// ReadFloat64 read float64
func (it *Iterator) ReadFloat64() (ret float64) {
	c := it.nextToken()
	if c == '-' {
		return -it.readPositiveFloat64()
	}
	it.unreadByte()
	return it.readPositiveFloat64()
}

func (it *Iterator) readPositiveFloat64() (ret float64) {
	i := it.head
	// First char.
	if i == it.tail {
		return it.readFloat64SlowPath()
	}
	c := it.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return it.readFloat64SlowPath()
	case endOfNumber:
		it.ReportError("readFloat64", "empty number")
		return
	case dotInNumber:
		it.ReportError("readFloat64", "leading dot is invalid")
		return
	case 0:
		if i == it.tail {
			return it.readFloat64SlowPath()
		}
		c = it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			it.ReportError("readFloat64", "leading zero is invalid")
			return
		}
	}
	value := uint64(ind)
	// Chars before dot.
NonDecimal:
	for ; i < it.tail; i++ {
		c = it.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return it.readFloat64SlowPath()
		case endOfNumber:
			it.head = i
			return float64(value)
		case dotInNumber:
			break NonDecimal
		}
		if value > uint64SafeToMultiple10 {
			return it.readFloat64SlowPath()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == it.tail {
			return it.readFloat64SlowPath()
		}
		for ; i < it.tail; i++ {
			c = it.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					it.head = i
					return float64(value) / float64(pow10[decimalPlaces])
				}
				// too many decimal places
				return it.readFloat64SlowPath()
			case invalidCharForNumber, dotInNumber:
				return it.readFloat64SlowPath()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return it.readFloat64SlowPath()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
			if value > maxFloat64 {
				return it.readFloat64SlowPath()
			}
		}
	}
	return it.readFloat64SlowPath()
}

func (it *Iterator) readFloat64SlowPath() (ret float64) {
	str := it.readNumberAsString()
	if it.Error != nil && it.Error != io.EOF {
		return
	}
	errMsg := validateFloat(str)
	if errMsg != "" {
		it.ReportError("readFloat64SlowPath", errMsg)
		return
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		it.Error = err
		return
	}
	return val
}

func validateFloat(str string) string {
	// strconv.ParseFloat is not validating `1.` or `1.e1`
	if len(str) == 0 {
		return "empty number"
	}
	if str[0] == '-' {
		return "-- is not valid"
	}
	dotPos := strings.IndexByte(str, tDot)
	if dotPos != -1 {
		if dotPos == len(str)-1 {
			return "dot can not be last character"
		}
		switch str[dotPos+1] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return "missing digit after dot"
		}
	}
	return ""
}

// ReadNumber read jir.Number
func (it *Iterator) ReadNumber() (ret json.Number) {
	return json.Number(it.readNumberAsString())
}
