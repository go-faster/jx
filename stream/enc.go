package stream

import "io"

// Encoder writes json tokens to given writer.
type Encoder[W io.Writer] struct {
	w      writer[W]
	indent int // count of spaces for single indentation level

	// first handles state for comma and indentation writing.
	//
	// New Object or Array appends new level to this slice, and
	// last element of this slice denotes whether first element was written.
	//
	// We write commas only before non-first element of Array or Object.
	//
	// See comma, begin, end and FieldStart for implementation details.
	//
	// Note: probably, this can be optimized as bit set to ease memory
	// consumption.
	//
	// See https://yourbasic.org/algorithms/your-basic-int/#simple-sets
	first []bool
}

// NewEncoder returns new Encoder that writes to given writer.
func NewEncoder[W io.Writer](w W) *Encoder[W] {
	return &Encoder[W]{
		w: writer[W]{
			writer: w,
			// TODO(tdakkota): add option to set buffer/size.
			buf: make([]byte, 0, 512),
		},
	}
}

// Close flushes all buffered data to underlying writer.
func (e *Encoder[W]) Close() error {
	e.w.flush()
	return e.w.flushErr
}

// SetIdent sets length of single indentation step.
func (e *Encoder[W]) SetIdent(n int) {
	e.indent = n
}

// Reset resets underlying buffer.
func (e *Encoder[W]) Reset(w W) {
	e.w.Reset(w)
	e.first = e.first[:0]
}

// RawStr writes string as raw json.
func (e *Encoder[W]) RawStr(v string) bool {
	return e.comma() || e.w.writeString(v)
}

// Raw writes byte slice as raw json.
func (e *Encoder[W]) Raw(b []byte) bool {
	return e.comma() || e.w.writeBytes(b...)
}

func (e *Encoder[W]) writeIndent() bool {
	if e.indent == 0 {
		return false
	}
	if e.w.writeByte('\n') {
		return true
	}
	for i := 0; i < len(e.first)*e.indent; i++ {
		if e.w.writeByte(' ') {
			return true
		}
	}
	return false
}
