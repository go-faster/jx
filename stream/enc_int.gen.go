package stream

import "io"

// FIXME(tdakkota): this is a straight copy-paste from w_int.gen.go
// 	Generator should be modified to generate this code.

var digits []uint32 // FIXME(tdakkota): re-use set

func init() {
	digits = make([]uint32, 1000)
	for i := uint32(0); i < 1000; i++ {
		digits[i] = (((i / 100) + '0') << 16) + ((((i / 10) % 10) + '0') << 8) + i%10 + '0'
		if i < 10 {
			digits[i] += 2 << 24
		} else if i < 100 {
			digits[i] += 1 << 24
		}
	}
}

func writeFirstBuf[W io.Writer](w *writer[W], v uint32) bool {
	r := make([]byte, 0, 4)
	start := v >> 24
	if start == 0 {
		r = append(r, byte(v>>16), byte(v>>8))
	} else if start == 1 {
		r = append(r, byte(v>>8))
	}
	r = append(r, byte(v))
	return w.writeBytes(r...)
}

func writeBuf[W io.Writer](w *writer[W], v uint32) bool {
	return w.writeBytes(byte(v>>16), byte(v>>8), byte(v))
}

// writeUInt16 encodes uint16.
func (e *Encoder[W]) writeUInt16(v uint16) (ok bool) {
	q0 := v
	// Iteration 0.
	q1 := q0 / 1000
	if q1 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q0])
		return ok
	}
	// Iteration 1.
	r1 := q0 - q1*1000
	ok = ok || writeFirstBuf(&e.w, digits[q1])
	ok = ok || writeBuf(&e.w, digits[r1])
	return ok
}

// UInt16 encodes uint16.
func (e *Encoder[W]) UInt16(v uint16) bool {
	return e.comma() || e.writeUInt16(v)
}

// Int16 encodes int16.
func (e *Encoder[W]) writeInt16(v int16) bool {
	var val uint16
	if v < 0 {
		val = uint16(-v)
		if e.w.writeByte('-') {
			return true
		}
	} else {
		val = uint16(v)
	}
	return e.writeUInt16(val)
}

// Int16 encodes int16.
func (e *Encoder[W]) Int16(v int16) bool {
	return e.comma() || e.writeInt16(v)
}

// writeUInt32 encodes uint32.
func (e *Encoder[W]) writeUInt32(v uint32) (ok bool) {
	q0 := v
	// Iteration 0.
	q1 := q0 / 1000
	if q1 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q0])
		return ok
	}
	// Iteration 1.
	r1 := q0 - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q1])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 2.
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q2])
		ok = ok || writeBuf(&e.w, digits[r2])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 3.
	r3 := q2 - q3*1000
	ok = ok || writeFirstBuf(&e.w, digits[q3])
	ok = ok || writeBuf(&e.w, digits[r3])
	ok = ok || writeBuf(&e.w, digits[r2])
	ok = ok || writeBuf(&e.w, digits[r1])
	return ok
}

// UInt32 encodes uint32.
func (e *Encoder[W]) UInt32(v uint32) bool {
	return e.comma() || e.writeUInt32(v)
}

// Int32 encodes int32.
func (e *Encoder[W]) writeInt32(v int32) bool {
	var val uint32
	if v < 0 {
		val = uint32(-v)
		if e.w.writeByte('-') {
			return true
		}
	} else {
		val = uint32(v)
	}
	return e.writeUInt32(val)
}

// Int32 encodes int32.
func (e *Encoder[W]) Int32(v int32) bool {
	return e.comma() || e.writeInt32(v)
}

// writeUInt64 encodes uint64.
func (e *Encoder[W]) writeUInt64(v uint64) (ok bool) {
	q0 := v
	// Iteration 0.
	q1 := q0 / 1000
	if q1 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q0])
		return ok
	}
	// Iteration 1.
	r1 := q0 - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q1])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 2.
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q2])
		ok = ok || writeBuf(&e.w, digits[r2])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 3.
	r3 := q2 - q3*1000
	q4 := q3 / 1000
	if q4 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q3])
		ok = ok || writeBuf(&e.w, digits[r3])
		ok = ok || writeBuf(&e.w, digits[r2])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 4.
	r4 := q3 - q4*1000
	q5 := q4 / 1000
	if q5 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q4])
		ok = ok || writeBuf(&e.w, digits[r4])
		ok = ok || writeBuf(&e.w, digits[r3])
		ok = ok || writeBuf(&e.w, digits[r2])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 5.
	r5 := q4 - q5*1000
	q6 := q5 / 1000
	if q6 == 0 {
		ok = ok || writeFirstBuf(&e.w, digits[q5])
		ok = ok || writeBuf(&e.w, digits[r5])
		ok = ok || writeBuf(&e.w, digits[r4])
		ok = ok || writeBuf(&e.w, digits[r3])
		ok = ok || writeBuf(&e.w, digits[r2])
		ok = ok || writeBuf(&e.w, digits[r1])
		return ok
	}
	// Iteration 6.
	r6 := q5 - q6*1000
	ok = ok || writeFirstBuf(&e.w, digits[q6])
	ok = ok || writeBuf(&e.w, digits[r6])
	ok = ok || writeBuf(&e.w, digits[r5])
	ok = ok || writeBuf(&e.w, digits[r4])
	ok = ok || writeBuf(&e.w, digits[r3])
	ok = ok || writeBuf(&e.w, digits[r2])
	ok = ok || writeBuf(&e.w, digits[r1])
	return ok
}

// UInt64 encodes uint64.
func (e *Encoder[W]) UInt64(v uint64) bool {
	return e.comma() || e.writeUInt64(v)
}

// Int64 encodes int64.
func (e *Encoder[W]) writeInt64(v int64) bool {
	var val uint64
	if v < 0 {
		val = uint64(-v)
		if e.w.writeByte('-') {
			return true
		}
	} else {
		val = uint64(v)
	}
	return e.writeUInt64(val)
}

// Int64 encodes int64.
func (e *Encoder[W]) Int64(v int64) bool {
	return e.comma() || e.writeInt64(v)
}
