package jx

import (
	"bytes"
	"io"
	"math/big"
	"strconv"

	"github.com/go-faster/errors"
)

var (
	pow10       = [...]uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
	floatDigits [256]int8
)

const (
	invalidCharForNumber = int8(-1)
	endOfNumber          = int8(-2)
	dotInNumber          = int8(-3)
)

func init() {
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
	str, err := d.numberAppend(nil)
	if err != nil {
		return nil, errors.Wrap(err, "number")
	}
	prec := 64
	if len(str) > prec {
		prec = len(str)
	}
	val, _, err := big.ParseFloat(string(str), 10, uint(prec), big.ToZero)
	if err != nil {
		return nil, errors.Wrap(err, "float")
	}
	return val, nil
}

// BigInt read big.Int
func (d *Decoder) BigInt() (*big.Int, error) {
	str, err := d.numberAppend(nil)
	if err != nil {
		return nil, errors.Wrap(err, "number")
	}
	v := big.NewInt(0)
	var ok bool
	if v, ok = v.SetString(string(str), 10); !ok {
		return nil, errors.New("invalid")
	}
	return v, nil
}

// Float32 reads float32 value.
func (d *Decoder) Float32() (float32, error) {
	c, err := d.more()
	if err != nil {
		return 0, errors.Wrap(err, "byte")
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
		return 0, errors.New("empty")
	case dotInNumber:
		return 0, errors.New("leading dot")
	case 0:
		if i == d.tail {
			return d.f32Slow()
		}
		c = d.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, errors.New("leading zero")
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

func (d *Decoder) number() []byte {
	start := d.head
	buf := d.buf[d.head:d.tail]
	for i, c := range buf {
		switch c {
		case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		default:
			// End of number.
			d.head += i
			return d.buf[start:d.head]
		}
	}
	// Buffer is number within head:tail.
	d.head = d.tail
	return d.buf[start:d.tail]
}

func (d *Decoder) numberAppend(b []byte) ([]byte, error) {
	for {
		b = append(b, d.number()...)
		if d.head != d.tail {
			return b, nil
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
	c, err := d.more()
	if err != nil {
		return 0, errors.Wrap(err, "byte")
	}

	switch c {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.atof64(c)
	default:
		return 0, badToken(c)
	}
}

func (d *Decoder) floatSlow(size int) (float64, error) {
	var buf [32]byte

	str, err := d.numberAppend(buf[:0])
	if err != nil {
		return 0, errors.Wrap(err, "number")
	}
	if err := validateFloat(str); err != nil {
		return 0, errors.Wrap(err, "invalid")
	}

	val, err := strconv.ParseFloat(string(str), size)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func validateFloat(str []byte) error {
	// strconv.ParseFloat is not validating `1.` or `1.e1`
	if len(str) == 0 {
		return errors.New("empty")
	}
	if str[0] == '-' {
		return errors.New("double minus")
	}
	if len(str) >= 2 && str[0] == '0' {
		switch str[1] {
		case 'e', 'E', '.':
		default:
			return errors.New("leading zero")
		}
	}
	dotPos := bytes.IndexByte(str, '.')
	if dotPos != -1 {
		if dotPos == len(str)-1 {
			return errors.New("dot as last char")
		}
		switch str[dotPos+1] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return errors.New("no digit after dot")
		}
	}
	return nil
}
