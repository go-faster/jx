package jx

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"strconv"

	"golang.org/x/xerrors"
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

// BigFloat read big.Float
func (it *Iterator) BigFloat() (*big.Float, error) {
	str, err := it.number(nil)
	if err != nil {
		return nil, xerrors.Errorf("number: %w", err)
	}
	prec := 64
	if len(str) > prec {
		prec = len(str)
	}
	val, _, err := big.ParseFloat(string(str), 10, uint(prec), big.ToZero)
	if err != nil {
		return nil, xerrors.Errorf("float: %w", err)
	}
	return val, nil
}

// BigInt read big.Int
func (it *Iterator) BigInt() (*big.Int, error) {
	str, err := it.number(nil)
	if err != nil {
		return nil, xerrors.Errorf("number: %w", err)
	}
	v := big.NewInt(0)
	var ok bool
	if v, ok = v.SetString(string(str), 10); !ok {
		return nil, xerrors.New("invalid")
	}
	return v, nil
}

// Float32 reads float32 value.
func (it *Iterator) Float32() (float32, error) {
	c, err := it.next()
	if err != nil {
		return 0, xerrors.Errorf("next: %w", err)
	}
	if c != '-' {
		it.unread()
	}
	v, err := it.positiveFloat32()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		v *= -1
	}
	return v, nil
}

func (it *Iterator) positiveFloat32() (float32, error) {
	i := it.head
	// First char.
	if i == it.tail {
		return it.f32Slow()
	}
	c := it.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return it.f32Slow()
	case endOfNumber:
		return 0, xerrors.New("empty")
	case dotInNumber:
		return 0, xerrors.New("leading dot")
	case 0:
		if i == it.tail {
			return it.f32Slow()
		}
		c = it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, xerrors.New("leading zero")
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
			return it.f32Slow()
		case endOfNumber:
			it.head = i
			return float32(value), nil
		case dotInNumber:
			break NonDecimalLoop
		}
		if value > uint64SafeToMultiple10 {
			return it.f32Slow()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == tDot {
		i++
		decimalPlaces := 0
		if i == it.tail {
			return it.f32Slow()
		}
		for ; i < it.tail; i++ {
			c = it.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					it.head = i
					return float32(float64(value) / float64(pow10[decimalPlaces])), nil
				}
				// too many decimal places
				return it.f32Slow()
			case invalidCharForNumber, dotInNumber:
				return it.f32Slow()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return it.f32Slow()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return it.f32Slow()
}

func (it *Iterator) number(b []byte) ([]byte, error) {
	for {
		for i := it.head; i < it.tail; i++ {
			switch c := it.buf[i]; c {
			case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				b = append(b, c)
				continue
			default:
				// End of number.
				it.head = i
				return b, nil
			}
		}
		if err := it.read(); err != nil {
			if err == io.EOF {
				return b, nil
			}
			return b, err
		}
	}
}

const (
	size32 = 32
	size64 = 64
)

func (it *Iterator) f32Slow() (float32, error) {
	v, err := it.floatSlow(size32)
	if err != nil {
		return 0, err
	}
	return float32(v), err
}

// Float64 read float64
func (it *Iterator) Float64() (float64, error) {
	c, err := it.next()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		v, err := it.positiveFloat64()
		if err != nil {
			return 0, err
		}
		return -v, err
	}
	it.unread()
	return it.positiveFloat64()
}

func (it *Iterator) positiveFloat64() (float64, error) {
	i := it.head
	// First char.
	if i == it.tail {
		return it.float64Slow()
	}
	c := it.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return it.float64Slow()
	case endOfNumber:
		return 0, xerrors.New("empty")
	case dotInNumber:
		return 0, xerrors.New("leading dot")
	case 0:
		if i == it.tail {
			return it.float64Slow()
		}
		c = it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, xerrors.New("leading zero")
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
			return it.float64Slow()
		case endOfNumber:
			it.head = i
			return float64(value), nil
		case dotInNumber:
			break NonDecimal
		}
		if value > uint64SafeToMultiple10 {
			return it.float64Slow()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == it.tail {
			return it.float64Slow()
		}
		for ; i < it.tail; i++ {
			c = it.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					it.head = i
					return float64(value) / float64(pow10[decimalPlaces]), nil
				}
				// too many decimal places
				return it.float64Slow()
			case invalidCharForNumber, dotInNumber:
				return it.float64Slow()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return it.float64Slow()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
			if value > maxFloat64 {
				return it.float64Slow()
			}
		}
	}
	return it.float64Slow()
}

func (it *Iterator) floatSlow(size int) (float64, error) {
	var buf [32]byte

	str, err := it.number(buf[:0])
	if err != nil {
		return 0, xerrors.Errorf("number: %w", err)
	}
	if err := validateFloat(str); err != nil {
		return 0, xerrors.Errorf("invalid: %w", err)
	}

	val, err := strconv.ParseFloat(string(str), size)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (it *Iterator) float64Slow() (float64, error) { return it.floatSlow(size64) }

func validateFloat(str []byte) error {
	// strconv.ParseFloat is not validating `1.` or `1.e1`
	if len(str) == 0 {
		return xerrors.New("empty")
	}
	if str[0] == '-' {
		return xerrors.New("double minus")
	}
	dotPos := bytes.IndexByte(str, tDot)
	if dotPos != -1 {
		if dotPos == len(str)-1 {
			return xerrors.New("dot as last char")
		}
		switch str[dotPos+1] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return xerrors.New("no digit after dot")
		}
	}
	return nil
}

// Number reads json.Number.
func (it *Iterator) Number() (json.Number, error) {
	str, err := it.number(nil)
	if err != nil {
		return "", err
	}
	return json.Number(str), nil
}
