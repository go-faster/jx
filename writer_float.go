package jx

import (
	"math"
	"strconv"

	"golang.org/x/xerrors"
)

var pow10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}

// WriteFloat32 write float32 to stream.
func (w *Writer) WriteFloat32(val float32) error {
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
	w.buf = strconv.AppendFloat(w.buf, float64(val), f, -1, 32)
	return nil
}

// WriteFloat32Lossy write float32 to stream with ONLY 6 digits precision although much much faster
func (w *Writer) WriteFloat32Lossy(val float32) error {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		return xerrors.Errorf("bad value %v", val)
	}
	if val < 0 {
		w.byte(tMinus)
		val = -val
	}
	if val > 0x4ffffff {
		return w.WriteFloat32(val)
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	w.Uint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return nil
	}
	w.byte(tDot)
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		w.byte(tZero)
	}
	w.Uint64(fval)
	for w.buf[len(w.buf)-1] == tZero {
		w.buf = w.buf[:len(w.buf)-1]
	}
	return nil
}

// WriteFloat64 write float64 to stream
func (w *Writer) WriteFloat64(val float64) error {
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
	start := len(w.buf)
	w.buf = strconv.AppendFloat(w.buf, val, f, -1, 64)
	if f == 'e' {
		return nil
	}

	// Ensure that we are still float.
	for _, c := range w.buf[start:] {
		if c == tDot {
			return nil
		}
	}
	w.buf = appendRune(w.buf, tDot)
	w.buf = appendRune(w.buf, tZero)
	return nil
}

// WriteFloat64Lossy write float64 to stream with ONLY 6 digits precision although much much faster
func (w *Writer) WriteFloat64Lossy(val float64) error {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return xerrors.Errorf("unsupported value: %f", val)
	}
	if val < 0 {
		w.byte(tMinus)
		val = -val
	}
	if val > 0x4ffffff {
		return w.WriteFloat64(val)
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	w.Uint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return nil
	}
	w.byte(tDot)
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		w.byte(tZero)
	}
	w.Uint64(fval)
	for w.buf[len(w.buf)-1] == tZero {
		w.buf = w.buf[:len(w.buf)-1]
	}
	return nil
}
