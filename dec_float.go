package jx

import (
	"io"
	"math/big"

	"github.com/go-faster/errors"
)

var (
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

type floatInfo struct {
	mantbits uint
	expbits  uint
	bias     int
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

	switch c {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.atof32(c)
	default:
		return 0, badToken(c)
	}
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

// readFloat reads a decimal mantissa and exponent from Decoder.
func (d *Decoder) readFloat(c byte) (mantissa uint64, exp int, neg, trunc bool, _ error) {
	const (
		digitTag  byte = 1
		closerTag byte = 2
	)
	// digits
	var (
		maxMantDigits = 19 // 10^19 fits in uint64
		nd            = 0
		ndMant        = 0
		dp            = 0

		sawDot = false

		e     = 0
		eSign = 0
	)
	defer func() {
		if !sawDot {
			dp = nd
		}

		if eSign != 0 {
			dp += e * eSign
		}
		if mantissa != 0 {
			exp = dp - ndMant
		}
	}()

	// Check that buffer is not empty.
	switch c {
	case '-':
		neg = true

		c, err := d.byte()
		if err != nil {
			return 0, 0, false, false, err
		}
		// Character after '-' must be a digit.
		if skipNumberSet[c] != digitTag {
			return 0, 0, false, false, badToken(c)
		}
		if c != '0' {
			d.unread()
			break
		}
		fallthrough
	case '0':
		dp--

		// If buffer is empty, try to read more.
		if d.head == d.tail {
			err := d.read()
			if err != nil {
				// There is no data anymore.
				if err == io.EOF {
					return
				}
				return 0, 0, false, false, err
			}
		}

		c = d.buf[d.head]
		if skipNumberSet[c] == closerTag {
			return
		}
		switch c {
		case '.':
			goto stateDot
		case 'e', 'E':
			goto stateExp
		default:
			return 0, 0, false, false, badToken(c)
		}
	default:
		d.unread()
	}
	for {
		for i, c := range d.buf[d.head:d.tail] {
			switch skipNumberSet[c] {
			case closerTag:
				d.head += i
				return
			case digitTag:
				nd++
				if ndMant < maxMantDigits {
					mantissa = (mantissa << 3) + (mantissa << 1) + uint64(floatDigits[c])
					ndMant++
				} else if c != '0' {
					trunc = true
				}
				continue
			}

			switch c {
			case '.':
				d.head += i
				goto stateDot
			case 'e', 'E':
				d.head += i
				goto stateExp
			default:
				return 0, 0, false, false, badToken(c)
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				d.head = d.tail
				return
			}
			return 0, 0, false, false, err
		}
	}

stateDot:
	dp = nd
	sawDot = true
	d.head++
	{
		var last byte = '.'
		for {
			for i, c := range d.buf[d.head:d.tail] {
				switch skipNumberSet[c] {
				case closerTag:
					d.head += i
					// Check that dot is not last character.
					if last == '.' {
						return 0, 0, false, false, io.ErrUnexpectedEOF
					}
					return
				case digitTag:
					last = c

					if c == '0' && nd == 0 {
						dp--
						continue
					}
					nd++
					if ndMant < maxMantDigits {
						mantissa = (mantissa << 3) + (mantissa << 1) + uint64(floatDigits[c])
						ndMant++
					} else if c != '0' {
						trunc = true
					}
					continue
				}

				switch c {
				case 'e', 'E':
					if last == '.' {
						return 0, 0, false, false, badToken(c)
					}
					d.head += i
					goto stateExp
				default:
					return 0, 0, false, false, badToken(c)
				}
			}

			if err := d.read(); err != nil {
				// There is no data anymore.
				if err == io.EOF {
					d.head = d.tail
					// Check that dot is not last character.
					if last == '.' {
						return 0, 0, false, false, io.ErrUnexpectedEOF
					}
					return
				}
				return 0, 0, false, false, err
			}
		}
	}
stateExp:
	d.head++
	eSign = 1
	// There must be a number or sign after e.
	{
		numOrSign, err := d.byte()
		if err != nil {
			return 0, 0, false, false, err
		}
		if skipNumberSet[numOrSign] != digitTag { // If next character is not a digit, check for sign.
			if numOrSign == '-' || numOrSign == '+' {
				if numOrSign == '-' {
					eSign = -1
				}
				num, err := d.byte()
				if err != nil {
					return 0, 0, false, false, err
				}
				// There must be a number after sign.
				if skipNumberSet[num] != digitTag {
					return 0, 0, false, false, badToken(num)
				}
				e = e*10 + int(num) - '0'
			} else {
				return 0, 0, false, false, badToken(numOrSign)
			}
		} else {
			e = e*10 + int(numOrSign) - '0'
		}
	}
	for {
		for i, c := range d.buf[d.head:d.tail] {
			if skipNumberSet[c] == closerTag {
				d.head += i
				return
			}
			if skipNumberSet[c] == 0 {
				return 0, 0, false, false, badToken(c)
			}
			if e < 10000 {
				e = e*10 + int(c) - '0'
			}
		}

		if err := d.read(); err != nil {
			// There is no data anymore.
			if err == io.EOF {
				d.head = d.tail
				return
			}
			return 0, 0, false, false, err
		}
	}
}
