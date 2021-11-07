package jx

import (
	"io"
)

// Encoder encodes json to underlying buffer.
//
// Zero value is valid.
type Encoder struct {
	buf    []byte // underlying buffer
	ident  int    // indentation step
	spaces int    // count of spaces

	// first handles state for commas writing.
	//
	// New Object or Array appends new level to this slice, and
	// last element of this slice denotes whether first element was written.
	//
	// We write commas only before non-first element of Array or Object.
	//
	// See comma, begin, end and Field for implementation details.
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
	e.ident = n
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

// Raw writes string as raw json.
func (e *Encoder) Raw(v string) {
	e.buf = append(e.buf, v...)
}

// RawBytes writes byte slice as raw json.
func (e *Encoder) RawBytes(b []byte) {
	e.buf = append(e.buf, b...)
}

// Null writes null to stream.
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

// Bool writes boolean.
func (e *Encoder) Bool(val bool) {
	if val {
		e.True()
	} else {
		e.False()
	}
}

// ObjStart writes object start, performing indentation if needed
func (e *Encoder) ObjStart() {
	e.comma()
	e.spaces += e.ident
	e.byte('{')
	e.writeIdent(0)
	e.begin()
}

// Field writes field name and colon.
//
// For non-zero indentation also writes single space after colon.
func (e *Encoder) Field(field string) {
	e.Str(field)
	if e.spaces > 0 {
		e.twoBytes(':', ' ')
	} else {
		e.byte(':')
	}
	if len(e.first) > 0 {
		e.first[e.last()] = true
	}
}

// ObjEnd writes end of object token, performing indentation if needed.
func (e *Encoder) ObjEnd() {
	e.writeIdent(e.ident)
	e.spaces -= e.ident
	e.byte('}')
	e.end()
}

// ObjEmpty writes empty object.
func (e *Encoder) ObjEmpty() {
	e.comma()
	e.byte('{')
	e.byte('}')
}

// ArrStart writes start of array, performing indentation if needed.
func (e *Encoder) ArrStart() {
	e.comma()
	e.spaces += e.ident
	e.byte('[')
	e.writeIdent(0)
	e.begin()
}

// ArrEmpty writes empty array.
func (e *Encoder) ArrEmpty() {
	e.comma()
	e.twoBytes('[', ']')
}

// ArrEnd writes end of array, performing indentation if needed.
func (e *Encoder) ArrEnd() {
	e.writeIdent(e.ident)
	e.spaces -= e.ident
	e.byte(']')
	e.end()
}

func (e *Encoder) writeIdent(delta int) {
	if e.spaces == 0 {
		return
	}
	e.byte('\n')
	spaces := e.spaces - delta
	for i := 0; i < spaces; i++ {
		e.buf = append(e.buf, ' ')
	}
}
