package jx

//go:generate go run ./tools/mkencint -output enc_int.gen.go

// Int encodes int.
func (e *Encoder) Int(v int) {
	e.Int64(int64(v))
}

// Uint encodes uint.
func (e *Encoder) Uint(v uint) {
	e.Uint64(uint64(v))
}

// Uint8 encodes uint8.
func (e *Encoder) Uint8(v uint8) {
	e.comma()
	// v is always smaller than digits size (1000)
	e.buf = writeFirstBuf(e.buf, digits[v])
}

// Int8 encodes int8.
func (e *Encoder) Int8(v int8) {
	e.comma()
	var val uint8
	if v < 0 {
		val = uint8(-v)
		e.comma()
		e.resetComma()
		e.buf = append(e.buf, '-')
	} else {
		val = uint8(v)
	}
	e.Uint8(val)
}
