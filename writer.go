package jx

import (
	"io"
)

// Writer writes json to io.Writer.
type Writer struct {
	out io.Writer
	buf []byte

	ident    int
	curIdent int
}

// SetIdent sets length of single indentation step.
func (w *Writer) SetIdent(n int) {
	w.ident = n
}

// NewWriter create new writer instance.
// out can be nil if write to internal buffer.
// bufSize is the initial size for the internal buffer in bytes.
func NewWriter(out io.Writer, bufSize int) *Writer {
	return &Writer{
		out: out,
		buf: make([]byte, 0, bufSize),
	}
}

// Reset reuse this stream instance by assign a new writer
func (w *Writer) Reset(out io.Writer) {
	w.out = out
	w.buf = w.buf[:0]
}

// Buf returns underlying buffer.
func (w *Writer) Buf() []byte { return w.buf }

// SetBuf allows to set the internal buffer directly.
func (w *Writer) SetBuf(buf []byte) { w.buf = buf }

// byte writes a single byte.
func (w *Writer) byte(c byte) {
	w.buf = append(w.buf, c)
}

func (w *Writer) twoBytes(c1, c2 byte) {
	w.buf = append(w.buf, c1, c2)
}

func (w *Writer) threeBytes(c1, c2, c3 byte) {
	w.buf = append(w.buf, c1, c2, c3)
}

func (w *Writer) fourBytes(c1, c2, c3, c4 byte) {
	w.buf = append(w.buf, c1, c2, c3, c4)
}

func (w *Writer) fiveBytes(c1, c2, c3, c4, c5 byte) {
	w.buf = append(w.buf, c1, c2, c3, c4, c5)
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	if w.out == nil {
		return nil
	}
	if _, err := w.out.Write(w.buf); err != nil {
		return err
	}
	w.buf = w.buf[:0]
	return nil
}

// Raw write string out without quotes, just like []byte.
func (w *Writer) Raw(v string) {
	w.buf = append(w.buf, v...)
}

// Null write null to stream.
func (w *Writer) Null() {
	w.fourBytes('n', 'u', 'l', 'l')
}

// True write true to stream.
func (w *Writer) True() {
	w.fourBytes('t', 'r', 'u', 'e')
}

// False writes false to stream.
func (w *Writer) False() {
	w.fiveBytes('f', 'a', 'l', 's', 'e')
}

// Bool writes boolean.
func (w *Writer) Bool(val bool) {
	if val {
		w.True()
	} else {
		w.False()
	}
}

// ObjStart writes { with possible indention.
func (w *Writer) ObjStart() {
	w.curIdent += w.ident
	w.byte('{')
	w.writeIdent(0)
}

// ObjField write "field": with possible indention.
func (w *Writer) ObjField(field string) {
	w.Str(field)
	if w.curIdent > 0 {
		w.twoBytes(':', ' ')
	} else {
		w.byte(':')
	}
}

// ObjEnd write } with possible indention
func (w *Writer) ObjEnd() {
	w.writeIdent(w.ident)
	w.curIdent -= w.ident
	w.byte('}')
}

// ObjEmpty write {}
func (w *Writer) ObjEmpty() {
	w.byte('{')
	w.byte('}')
}

// More write , with possible indention
func (w *Writer) More() {
	w.byte(',')
	w.writeIdent(0)
}

// ArrStart writes [ with possible indention.
func (w *Writer) ArrStart() {
	w.curIdent += w.ident
	w.byte('[')
	w.writeIdent(0)
}

// ArrEmpty writes [].
func (w *Writer) ArrEmpty() {
	w.twoBytes('[', ']')
}

// ArrEnd writes ] with possible indention.
func (w *Writer) ArrEnd() {
	w.writeIdent(w.ident)
	w.curIdent -= w.ident
	w.byte(']')
}

func (w *Writer) writeIdent(delta int) {
	if w.curIdent == 0 {
		return
	}
	w.byte('\n')
	toWrite := w.curIdent - delta
	for i := 0; i < toWrite; i++ {
		w.buf = append(w.buf, ' ')
	}
}
