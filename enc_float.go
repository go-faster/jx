package jx

// Float32 writes float32.
//
// NB: Infinities and NaN are represented as null.
func (e *Encoder) Float32(v float32) { e.float(float64(v), 32) }

// Float64 writes float64.
//
// NB: Infinities and NaN are represented as null.
func (e *Encoder) Float64(v float64) { e.float(v, 64) }
