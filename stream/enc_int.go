package stream

// Int encodes int.
func (e *Encoder[W]) Int(v int) bool {
	return e.Int64(int64(v))
}

// UInt encodes uint.
func (e *Encoder[W]) UInt(v uint) bool {
	return e.UInt64(uint64(v))
}

// UInt8 encodes uint8.
func (e *Encoder[W]) UInt8(v uint8) bool {
	// v is always smaller than digits size (1000)
	return writeFirstBuf(&e.w, digits[v])
}

// Int8 encodes int8.
func (e *Encoder[W]) Int8(v int8) bool {
	var val uint8
	if v < 0 {
		val = uint8(-v)
		if e.w.writeByte('-') {
			return true
		}
	} else {
		val = uint8(v)
	}
	return e.UInt8(val)
}
