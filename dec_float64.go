// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jx

var (
	float64info = floatInfo{52, 11, -1023}
	// Exact powers of 10.
	float64pow10 = [...]float64{
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
		1e20, 1e21, 1e22,
	}
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
