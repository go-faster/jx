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
func (e *Encoder) String() string {
	return string(e.Bytes())
}

// Reset resets underlying buffer.
func (e *Encoder) Reset() {
	e.buf = e.buf[:0]
}

// Bytes returns underlying buffer.
func (e *Encoder) Bytes() []byte { return e.buf }

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
	e.fourBytes('n', 'u', 'l', 'l')
}

// True writes true.
func (e *Encoder) True() {
	e.fourBytes('t', 'r', 'u', 'e')
}

// False writes false.
func (e *Encoder) False() {
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
	e.spaces += e.ident
	e.byte('{')
	e.writeIdent(0)
}

// ObjField writes field name and colon.
//
// For non-zero indentation also writes single space after colon.
func (e *Encoder) ObjField(field string) {
	e.Str(field)
	if e.spaces > 0 {
		e.twoBytes(':', ' ')
	} else {
		e.byte(':')
	}
}

// ObjEnd writes end of object token, performing indentation if needed.
func (e *Encoder) ObjEnd() {
	e.writeIdent(e.ident)
	e.spaces -= e.ident
	e.byte('}')
}

// ObjEmpty writes empty object.
func (e *Encoder) ObjEmpty() {
	e.byte('{')
	e.byte('}')
}

// More writes comma, performing indentation if needed.
func (e *Encoder) More() {
	e.byte(',')
	e.writeIdent(0)
}

// ArrStart writes start of array, performing indentation if needed.
func (e *Encoder) ArrStart() {
	e.spaces += e.ident
	e.byte('[')
	e.writeIdent(0)
}

// ArrEmpty writes empty array.
func (e *Encoder) ArrEmpty() {
	e.twoBytes('[', ']')
}

// ArrEnd writes end of array, performing indentation if needed.
func (e *Encoder) ArrEnd() {
	e.writeIdent(e.ident)
	e.spaces -= e.ident
	e.byte(']')
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
