package jx

import "io"

// Encoder encodes json to underlying buffer.
//
// Zero value is valid.
type Encoder struct {
	buf    []byte // underlying buffer
	indent int    // count of spaces for single indentation level

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

// Write implements io.Writer.
func (e *Encoder) Write(p []byte) (n int, err error) {
	e.buf = append(e.buf, p...)
	return len(p), nil
}

// WriteTo implements io.WriterTo.
func (e *Encoder) WriteTo(w io.Writer) (n int64, err error) {
	wrote, err := w.Write(e.buf)
	return int64(wrote), err
}

// SetIdent sets length of single indentation step.
func (e *Encoder) SetIdent(n int) {
	e.indent = n
}

// String returns string of underlying buffer.
func (e Encoder) String() string {
	return string(e.Bytes())
}

// Reset resets underlying buffer.
func (e *Encoder) Reset() {
	e.buf = e.buf[:0]
	e.first = e.first[:0]
}

// Bytes returns underlying buffer.
func (e Encoder) Bytes() []byte { return e.buf }

// SetBytes sets underlying buffer.
func (e *Encoder) SetBytes(buf []byte) { e.buf = buf }

// byte writes a single byte.
func (e *Encoder) byte(c byte) {
	e.buf = append(e.buf, c)
}

func (e *Encoder) twoBytes(c1, c2 byte) {
	e.buf = append(e.buf, c1, c2)
}

func (e *Encoder) threeBytes(c1, c2, c3 byte) {
	e.buf = append(e.buf, c1, c2, c3)
}

func (e *Encoder) fourBytes(c1, c2, c3, c4 byte) {
	e.buf = append(e.buf, c1, c2, c3, c4)
}

func (e *Encoder) fiveBytes(c1, c2, c3, c4, c5 byte) {
	e.buf = append(e.buf, c1, c2, c3, c4, c5)
}

// RawStr writes string as raw json.
func (e *Encoder) RawStr(v string) {
	e.buf = append(e.buf, v...)
}

// Raw writes byte slice as raw json.
func (e *Encoder) Raw(b []byte) {
	e.buf = append(e.buf, b...)
}

// Null writes null.
func (e *Encoder) Null() {
	e.comma()
	e.fourBytes('n', 'u', 'l', 'l')
}

// True writes true.
func (e *Encoder) True() {
	e.comma()
	e.fourBytes('t', 'r', 'u', 'e')
}

// False writes false.
func (e *Encoder) False() {
	e.comma()
	e.fiveBytes('f', 'a', 'l', 's', 'e')
}

// Bool encodes boolean.
func (e *Encoder) Bool(v bool) {
	if v {
		e.True()
	} else {
		e.False()
	}
}

// ObjStart writes object start, performing indentation if needed.
//
// Use Obj as convenience helper for writing objects.
func (e *Encoder) ObjStart() {
	e.comma()
	e.byte('{')
	e.begin()
	e.writeIndent()
}

// FieldStart encodes field name and writes colon.
//
// For non-zero indentation also writes single space after colon.
//
// Use Field as convenience helper for encoding fields.
func (e *Encoder) FieldStart(field string) {
	e.Str(field)
	if e.indent > 0 {
		e.twoBytes(':', ' ')
	} else {
		e.byte(':')
	}
	if len(e.first) > 0 {
		e.first[e.last()] = true
	}
}

// Field encodes field start and then invokes callback.
//
// Has ~5ns overhead over FieldStart.
func (e *Encoder) Field(name string, f func(e *Encoder)) {
	e.FieldStart(name)
	f(e)
}

// ObjEnd writes end of object token, performing indentation if needed.
//
// Use Obj as convenience helper for writing objects.
func (e *Encoder) ObjEnd() {
	e.end()
	e.writeIndent()
	e.byte('}')
}

// ObjEmpty writes empty object.
func (e *Encoder) ObjEmpty() {
	e.comma()
	e.twoBytes('{', '}')
}

// Obj writes start of object, invokes callback and writes end of object.
//
// If callback is nil, writes empty object.
func (e *Encoder) Obj(f func(e *Encoder)) {
	if f == nil {
		e.ObjEmpty()
		return
	}
	e.ObjStart()
	f(e)
	e.ObjEnd()
}

// ArrStart writes start of array, performing indentation if needed.
//
// Use Arr as convenience helper for writing arrays.
func (e *Encoder) ArrStart() {
	e.comma()
	e.byte('[')
	e.begin()
	e.writeIndent()
}

// ArrEmpty writes empty array.
func (e *Encoder) ArrEmpty() {
	e.comma()
	e.twoBytes('[', ']')
}

// ArrEnd writes end of array, performing indentation if needed.
//
// Use Arr as convenience helper for writing arrays.
func (e *Encoder) ArrEnd() {
	e.end()
	e.writeIndent()
	e.byte(']')
}

// Arr writes start of array, invokes callback and writes end of array.
//
// If callback is nil, writes empty array.
func (e *Encoder) Arr(f func(e *Encoder)) {
	if f == nil {
		e.ArrEmpty()
		return
	}
	e.ArrStart()
	f(e)
	e.ArrEnd()
}

func (e *Encoder) writeIndent() {
	if e.indent == 0 {
		return
	}
	e.byte('\n')
	for i := 0; i < len(e.first)*e.indent; i++ {
		e.buf = append(e.buf, ' ')
	}
}
