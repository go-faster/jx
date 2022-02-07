// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jx

import "io"

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

type floatInfo struct {
	mantbits uint
	expbits  uint
	bias     int
}

var (
	float32info = floatInfo{23, 8, -127}
	float64info = floatInfo{52, 11, -1023}

	// Exact powers of 10.
	float64pow10 = [...]float64{
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
		1e20, 1e21, 1e22,
	}
	float32pow10 = [...]float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
)

// If possible to convert decimal representation to 64-bit float f exactly,
// entirely in floating-point math, do so, avoiding the expense of decimalToFloatBits.
// Three common cases:
//	value is exact integer
//	value is exact integer * exact power of ten
//	value is exact integer / exact power of ten
// These all produce potentially inexact but correctly rounded answers.
func atof64exact(mantissa uint64, exp int, neg bool) (f float64, ok bool) {
	if mantissa>>float64info.mantbits != 0 {
		return
	}
	f = float64(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		// an integer.
		return f, true
	// Exact integers are <= 10^15.
	// Exact powers of ten are <= 10^22.
	case exp > 0 && exp <= 15+22: // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 22 {
			f *= float64pow10[exp-22]
			exp = 22
		}
		if f > 1e15 || f < -1e15 {
			// the exponent was really too large.
			return
		}
		return f * float64pow10[exp], true
	case exp < 0 && exp >= -22: // int / 10^k
		return f / float64pow10[-exp], true
	}
	return
}

// If possible to compute mantissa*10^exp to 32-bit float f exactly,
// entirely in floating-point math, do so, avoiding the machinery above.
func atof32exact(mantissa uint64, exp int, neg bool) (f float32, ok bool) {
	if mantissa>>float32info.mantbits != 0 {
		return
	}
	f = float32(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		return f, true
	// Exact integers are <= 10^7.
	// Exact powers of ten are <= 10^10.
	case exp > 0 && exp <= 7+10: // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 10 {
			f *= float32pow10[exp-10]
			exp = 10
		}
		if f > 1e7 || f < -1e7 {
			// the exponent was really too large.
			return
		}
		return f * float32pow10[exp], true
	case exp < 0 && exp >= -10: // int / 10^k
		return f / float32pow10[-exp], true
	}
	return
}

func (d *Decoder) atof64(c byte) (_ float64, err error) {
	mantissa, exp, neg, trunc, err := d.readFloat(c)
	if err != nil {
		return 0, err
	}

	// Try pure floating-point arithmetic conversion, and if that fails,
	// the Eisel-Lemire algorithm.
	if !trunc {
		if f, ok := atof64exact(mantissa, exp, neg); ok {
			return f, nil
		}
	}

	if f, ok := eiselLemire64(mantissa, exp, neg); ok {
		if !trunc {
			return f, nil
		}
		// Even if the mantissa was truncated, we may
		// have found the correct result. Confirm by
		// converting the upper mantissa bound.
		fUp, ok := eiselLemire64(mantissa+1, exp, neg)
		if ok && f == fUp {
			return f, nil
		}
	}

	return 0, err
}
