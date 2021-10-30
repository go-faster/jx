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
	floatDigits['.'] = dotInNumber
}

// BigFloat read big.Float
func (d *Decoder) BigFloat() (*big.Float, error) {
	str, err := d.number(nil)
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
func (d *Decoder) BigInt() (*big.Int, error) {
	str, err := d.number(nil)
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
func (d *Decoder) Float32() (float32, error) {
	c, err := d.next()
	if err != nil {
		return 0, xerrors.Errorf("next: %w", err)
	}
	if c != '-' {
		d.unread()
	}
	v, err := d.positiveFloat32()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		v *= -1
	}
	return v, nil
}

func (d *Decoder) positiveFloat32() (float32, error) {
	i := d.head
	// First char.
	if i == d.tail {
		return d.f32Slow()
	}
	c := d.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return d.f32Slow()
	case endOfNumber:
		return 0, xerrors.New("empty")
	case dotInNumber:
		return 0, xerrors.New("leading dot")
	case 0:
		if i == d.tail {
			return d.f32Slow()
		}
		c = d.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, xerrors.New("leading zero")
		}
	}
	value := uint64(ind)
	// Chars before dot.
NonDecimalLoop:
	for ; i < d.tail; i++ {
		c = d.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return d.f32Slow()
		case endOfNumber:
			d.head = i
			return float32(value), nil
		case dotInNumber:
			break NonDecimalLoop
		}
		if value > uint64SafeToMultiple10 {
			return d.f32Slow()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// Chars after dot.
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == d.tail {
			return d.f32Slow()
		}
		for ; i < d.tail; i++ {
			c = d.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					d.head = i
					return float32(float64(value) / float64(pow10[decimalPlaces])), nil
				}
				// too many decimal places
				return d.f32Slow()
			case invalidCharForNumber, dotInNumber:
				return d.f32Slow()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return d.f32Slow()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return d.f32Slow()
}

func (d *Decoder) number(b []byte) ([]byte, error) {
	for {
		for i := d.head; i < d.tail; i++ {
			switch c := d.buf[i]; c {
			case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				b = append(b, c)
				continue
			default:
				// End of number.
				d.head = i
				return b, nil
			}
		}
		if err := d.read(); err != nil {
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

func (d *Decoder) f32Slow() (float32, error) {
	v, err := d.floatSlow(size32)
	if err != nil {
		return 0, err
	}
	return float32(v), err
}

// Float64 read float64
func (d *Decoder) Float64() (float64, error) {
	c, err := d.next()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		v, err := d.positiveFloat64()
		if err != nil {
			return 0, err
		}
		return -v, err
	}
	d.unread()
	return d.positiveFloat64()
}

func (d *Decoder) positiveFloat64() (float64, error) {
	i := d.head
	// First char.
	if i == d.tail {
		return d.float64Slow()
	}
	c := d.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return d.float64Slow()
	case endOfNumber:
		return 0, xerrors.New("empty")
	case dotInNumber:
		return 0, xerrors.New("leading dot")
	case 0:
		if i == d.tail {
			return d.float64Slow()
		}
		c = d.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, xerrors.New("leading zero")
		}
	}
	value := uint64(ind)
	// Chars before dot.
NonDecimal:
	for ; i < d.tail; i++ {
		c = d.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return d.float64Slow()
		case endOfNumber:
			d.head = i
			return float64(value), nil
		case dotInNumber:
			break NonDecimal
		}
		if value > uint64SafeToMultiple10 {
			return d.float64Slow()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == d.tail {
			return d.float64Slow()
		}
		for ; i < d.tail; i++ {
			c = d.buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					d.head = i
					return float64(value) / float64(pow10[decimalPlaces]), nil
				}
				// too many decimal places
				return d.float64Slow()
			case invalidCharForNumber, dotInNumber:
				return d.float64Slow()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return d.float64Slow()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
			if value > maxFloat64 {
				return d.float64Slow()
			}
		}
	}
	return d.float64Slow()
}

func (d *Decoder) floatSlow(size int) (float64, error) {
	var buf [32]byte

	str, err := d.number(buf[:0])
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

func (d *Decoder) float64Slow() (float64, error) { return d.floatSlow(size64) }

func validateFloat(str []byte) error {
	// strconv.ParseFloat is not validating `1.` or `1.e1`
	if len(str) == 0 {
		return xerrors.New("empty")
	}
	if str[0] == '-' {
		return xerrors.New("double minus")
	}
	dotPos := bytes.IndexByte(str, '.')
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
func (d *Decoder) Number() (json.Number, error) {
	str, err := d.number(nil)
	if err != nil {
		return "", err
	}
	return json.Number(str), nil
}
