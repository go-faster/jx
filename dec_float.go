package jx

import (
	"bytes"
	"io"
	"math/big"

	"github.com/go-faster/errors"
)

var (
	pow10       = [...]uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
	floatDigits = [256]int8{}
)

const (
	invalidCharForNumber = int8(-1)
	endOfNumber          = int8(-2)
	dotInNumber          = int8(-3)
	maxFloat64           = 1<<63 - 1
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
		return d.float32Slow()
	}
	c := d.buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return d.float32Slow()
	case endOfNumber:
		return 0, errors.New("empty")
	case dotInNumber:
		return 0, errors.New("leading dot")
	case 0:
		if i == d.tail {
			return d.float32Slow()
		}
		c = d.buf[i]
		if floatDigits[c] >= 0 {
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
			return d.float32Slow()
		case endOfNumber:
			d.head = i
			return float32(value), nil
		case dotInNumber:
			break NonDecimalLoop
		}
		if value > uint64SafeToMultiple10 {
			return d.float32Slow()
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// Chars after dot.
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == d.tail {
			return d.float32Slow()
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
				return d.float32Slow()
			case invalidCharForNumber, dotInNumber:
				return d.float32Slow()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return d.float32Slow()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return d.float32Slow()
}

var numberSet = [256]byte{
	'+': 1,
	'-': 1,
	'.': 1,
	'e': 1,
	'E': 1,
	'0': 1,
	'1': 1,
	'2': 1,
	'3': 1,
	'4': 1,
	'5': 1,
	'6': 1,
	'7': 1,
	'8': 1,
	'9': 1,
}

func (d *Decoder) number() []byte {
	start := d.head
	buf := d.buf[d.head:d.tail]
	for i, c := range buf {
		if numberSet[c] == 0 {
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
	size64 = 64
)

func (d *Decoder) float32Slow() (float32, error) {
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
	if floatDigits[c] >= 0 {
		d.unread()
		return d.positiveFloat64()
	}
	switch c {
	case '-':
		v, err := d.positiveFloat64()
		if err != nil {
			return 0, err
		}
		return -v, err
	default:
		return 0, badToken(c)
	}
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
		return 0, errors.New("empty")
	case dotInNumber:
		return 0, errors.New("leading dot")
	case 0:
		if i == d.tail {
			return d.float64Slow()
		}
		c = d.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, errors.New("leading zero")
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
			// Not checking for uint64SafeToMultiple10 here because
			// if condition is positive value multiplied by 10 is
			// guaranteed to be bigger than maxFloat64.
			value = (value << 3) + (value << 1) + uint64(ind)
			if value > maxFloat64 {
				return d.float64Slow()
			}
		}
	}
	return d.float64Slow()
}

func (d *Decoder) float64Slow() (float64, error) { return d.floatSlow(size64) }

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
