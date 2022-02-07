package jx

var (
	float32info = floatInfo{23, 8, -127}
	// Exact powers of 10.
	float32pow10 = [...]float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
)

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

func (d *Decoder) atof32(c byte) (_ float32, err error) {
	mantissa, exp, neg, trunc, err := d.readFloat(c)
	if err != nil {
		return 0, err
	}

	// Try pure floating-point arithmetic conversion, and if that fails,
	// the Eisel-Lemire algorithm.
	if !trunc {
		if f, ok := atof32exact(mantissa, exp, neg); ok {
			return f, nil
		}
	}

	if f, ok := eiselLemire32(mantissa, exp, neg); ok {
		if !trunc {
			return f, nil
		}
		// Even if the mantissa was truncated, we may
		// have found the correct result. Confirm by
		// converting the upper mantissa bound.
		fUp, ok := eiselLemire32(mantissa+1, exp, neg)
		if ok && f == fUp {
			return f, nil
		}
	}

	return 0, err
}
