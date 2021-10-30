package jx

import (
	"math"
	"strconv"

	"golang.org/x/xerrors"
)

var pow10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}

// WriteFloat32 writes float32 to writer.
func (e *Encoder) WriteFloat32(val float32) error {
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

// Float32Lossy writes float32 to writer with ONLY 6 digits precision, but fast.
func (e *Encoder) Float32Lossy(val float32) error {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		return xerrors.Errorf("bad value %v", val)
	}
	if val < 0 {
		e.byte(tMinus)
		val = -val
	}
	if val > 0x4ffffff {
		return e.WriteFloat32(val)
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	e.Uint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return nil
	}
	e.byte(tDot)
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		e.byte(tZero)
	}
	e.Uint64(fval)
	for e.buf[len(e.buf)-1] == tZero {
		e.buf = e.buf[:len(e.buf)-1]
	}
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
		if c == tDot {
			return nil
		}
	}
	e.buf = appendRune(e.buf, tDot)
	e.buf = appendRune(e.buf, tZero)
	return nil
}

// Float64Lossy write float64 to stream with ONLY 6 digits precision although much much faster
func (e *Encoder) Float64Lossy(val float64) error {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return xerrors.Errorf("unsupported value: %f", val)
	}
	if val < 0 {
		e.byte(tMinus)
		val = -val
	}
	if val > 0x4ffffff {
		return e.Float64(val)
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	e.Uint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return nil
	}
	e.byte(tDot)
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		e.byte(tZero)
	}
	e.Uint64(fval)
	for e.buf[len(e.buf)-1] == tZero {
		e.buf = e.buf[:len(e.buf)-1]
	}
	return nil
}
