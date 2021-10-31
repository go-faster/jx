package jx

import (
	"math"
	"strconv"

	"golang.org/x/xerrors"
)

var pow10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}

// Float32 writes float32..
func (e *Encoder) Float32(val float32) error {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		return xerrors.Errorf("bad value %v", val)
	}
	abs := math.Abs(float64(val))
	f := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			f = 'e'
		}
	}
	e.buf = strconv.AppendFloat(e.buf, float64(val), f, -1, 32)
	return nil
}

// Float64 writes float64 to stream.
func (e *Encoder) Float64(val float64) error {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return xerrors.Errorf("unsupported value: %f", val)
	}
	abs := math.Abs(val)
	f := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			f = 'e'
		}
	}
	start := len(e.buf)
	e.buf = strconv.AppendFloat(e.buf, val, f, -1, 64)
	if f == 'e' {
		return nil
	}

	// Ensure that we are still float.
	for _, c := range e.buf[start:] {
		if c == '.' {
			return nil
		}
	}
	e.buf = appendRune(e.buf, '.')
	e.buf = appendRune(e.buf, '0')
	return nil
}
