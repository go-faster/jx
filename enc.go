package jx

import (
	"io"
)

// Encoder encodes json to internal buffer and flushes to io.Writer if set.
//
// Call Flush to write buffer to io.Writer.
type Encoder struct {
	out io.Writer
	buf []byte

	ident    int
	curIdent int
}

// SetIdent sets length of single indentation step.
func (e *Encoder) SetIdent(n int) {
	e.ident = n
}

// NewEncoder create new encoder instance.
func NewEncoder(out io.Writer, bufSize int) *Encoder {
	return &Encoder{
		out: out,
		buf: make([]byte, 0, bufSize),
	}
}

// Reset reuse this stream instance by assign a new writer
func (e *Encoder) Reset(out io.Writer) {
	e.out = out
	e.buf = e.buf[:0]
}

// Bytes returns underlying buffer.
func (e *Encoder) Bytes() []byte { return e.buf }

// SetBuf allows to set the internal buffer directly.
func (e *Encoder) SetBuf(buf []byte) { e.buf = buf }

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

// Flush writes any buffered data to the underlying io.Writer.
func (e *Encoder) Flush() error {
	if e.out == nil {
		return nil
	}
	if _, err := e.out.Write(e.buf); err != nil {
		return err
	}
	e.buf = e.buf[:0]
	return nil
}

// Raw writes raw json.
func (e *Encoder) Raw(v string) {
	e.buf = append(e.buf, v...)
}

// RawBytes writes raw json.
func (e *Encoder) RawBytes(b []byte) {
	e.buf = append(e.buf, b...)
}

// Null write null to stream.
func (e *Encoder) Null() {
	e.fourBytes('n', 'u', 'l', 'l')
}

// True write true to stream.
func (e *Encoder) True() {
	e.fourBytes('t', 'r', 'u', 'e')
}

// False writes false to stream.
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

// ObjStart writes { with possible indention.
func (e *Encoder) ObjStart() {
	e.curIdent += e.ident
	e.byte('{')
	e.writeIdent(0)
}

// ObjField write "field": with possible indention.
func (e *Encoder) ObjField(field string) {
	e.String(field)
	if e.curIdent > 0 {
		e.twoBytes(':', ' ')
	} else {
		e.byte(':')
	}
}

// ObjEnd write } with possible indention
func (e *Encoder) ObjEnd() {
	e.writeIdent(e.ident)
	e.curIdent -= e.ident
	e.byte('}')
}

// ObjEmpty write {}
func (e *Encoder) ObjEmpty() {
	e.byte('{')
	e.byte('}')
}

// More write , with possible indention
func (e *Encoder) More() {
	e.byte(',')
	e.writeIdent(0)
}

// ArrStart writes [ with possible indention.
func (e *Encoder) ArrStart() {
	e.curIdent += e.ident
	e.byte('[')
	e.writeIdent(0)
}

// ArrEmpty writes [].
func (e *Encoder) ArrEmpty() {
	e.twoBytes('[', ']')
}

// ArrEnd writes ] with possible indention.
func (e *Encoder) ArrEnd() {
	e.writeIdent(e.ident)
	e.curIdent -= e.ident
	e.byte(']')
}

func (e *Encoder) writeIdent(delta int) {
	if e.curIdent == 0 {
		return
	}
	e.byte('\n')
	toWrite := e.curIdent - delta
	for i := 0; i < toWrite; i++ {
		e.buf = append(e.buf, ' ')
	}
}
