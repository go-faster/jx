package jx

//go:generate go run ./tools/mkencint -output w_int.gen.go

// Int encodes int.
func (w *Writer) Int(v int) {
	w.Int64(int64(v))
}

// UInt encodes uint.
func (w *Writer) UInt(v uint) {
	w.UInt64(uint64(v))
}

// UInt8 encodes uint8.
func (w *Writer) UInt8(v uint8) {
	// v is always smaller than digits size (1000)
	w.Buf = writeFirstBuf(w.Buf, digits[v])
}

// Int8 encodes int8.
func (w *Writer) Int8(v int8) {
	var val uint8
	if v < 0 {
		val = uint8(-v)
		w.Buf = append(w.Buf, '-')
	} else {
		val = uint8(v)
	}
	w.UInt8(val)
}
