package jx

// Float32 encodes float32.
//
// NB: Infinities and NaN are represented as null.
func (w *Writer) Float32(v float32) { w.Float(float64(v), 32) }

// Float64 encodes float64.
//
// NB: Infinities and NaN are represented as null.
func (w *Writer) Float64(v float64) { w.Float(v, 64) }
