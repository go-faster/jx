package jir

import (
	"fmt"
	"math"
	"strconv"
)

var pow10 []uint64

func init() {
	pow10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
}

// WriteFloat32 write float32 to stream
func (s *Stream) WriteFloat32(val float32) {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		s.Error = fmt.Errorf("unsupported value: %f", val)
		return
	}
	abs := math.Abs(float64(val))
	f := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			f = 'e'
		}
	}
	s.buf = strconv.AppendFloat(s.buf, float64(val), f, -1, 32)
}

// WriteFloat32Lossy write float32 to stream with ONLY 6 digits precision although much much faster
func (s *Stream) WriteFloat32Lossy(val float32) {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		s.Error = fmt.Errorf("unsupported value: %f", val)
		return
	}
	if val < 0 {
		s.writeByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		s.WriteFloat32(val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	s.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	s.writeByte('.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		s.writeByte('0')
	}
	s.WriteUint64(fval)
	for s.buf[len(s.buf)-1] == '0' {
		s.buf = s.buf[:len(s.buf)-1]
	}
}

// WriteFloat64 write float64 to stream
func (s *Stream) WriteFloat64(val float64) {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		s.Error = fmt.Errorf("unsupported value: %f", val)
		return
	}
	abs := math.Abs(val)
	f := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			f = 'e'
		}
	}
	start := len(s.buf)
	s.buf = strconv.AppendFloat(s.buf, val, f, -1, 64)
	if f == 'e' {
		return
	}

	// Ensure that we are still float.
	for _, c := range s.buf[start:] {
		if c == '.' {
			return
		}
	}
	s.buf = appendRune(s.buf, '.')
	s.buf = appendRune(s.buf, '0')
}

// WriteFloat64Lossy write float64 to stream with ONLY 6 digits precision although much much faster
func (s *Stream) WriteFloat64Lossy(val float64) {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		s.Error = fmt.Errorf("unsupported value: %f", val)
		return
	}
	if val < 0 {
		s.writeByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		s.WriteFloat64(val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	s.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	s.writeByte('.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		s.writeByte('0')
	}
	s.WriteUint64(fval)
	for s.buf[len(s.buf)-1] == '0' {
		s.buf = s.buf[:len(s.buf)-1]
	}
}
