package jx

import "io"

// Writer writes json tokens to underlying buffer.
//
// Zero value is valid.
type Writer struct {
	Buf []byte // underlying buffer
}

// Write implements io.Writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.Buf = append(w.Buf, p...)
	return len(p), nil
}

// WriteTo implements io.WriterTo.
func (w *Writer) WriteTo(t io.Writer) (n int64, err error) {
	wrote, err := t.Write(w.Buf)
	return int64(wrote), err
}

// String returns string of underlying buffer.
func (w Writer) String() string {
	return string(w.Buf)
}

// Reset resets underlying buffer.
func (w *Writer) Reset() {
	w.Buf = w.Buf[:0]
}

// byte writes a single byte.
func (w *Writer) byte(c byte) {
	w.Buf = append(w.Buf, c)
}

func (w *Writer) twoBytes(c1, c2 byte) {
	w.Buf = append(w.Buf, c1, c2)
}

// RawStr writes string as raw json.
func (w *Writer) RawStr(v string) {
	w.rawStr(v)
}

func (w *Writer) rawStr(v string) {
	w.Buf = append(w.Buf, v...)
}

// Raw writes byte slice as raw json.
func (w *Writer) Raw(b []byte) {
	w.Buf = append(w.Buf, b...)
}

// Null writes null.
func (w *Writer) Null() {
	w.Buf = append(w.Buf, "null"...)
}

// True writes true.
func (w *Writer) True() {
	w.Buf = append(w.Buf, "true"...)
}

// False writes false.
func (w *Writer) False() {
	w.Buf = append(w.Buf, "false"...)
}

// Bool encodes boolean.
func (w *Writer) Bool(v bool) {
	if v {
		w.True()
	} else {
		w.False()
	}
}

// ObjStart writes object start.
func (w *Writer) ObjStart() {
	w.byte('{')
}

// FieldStart encodes field name and writes colon.
func (w *Writer) FieldStart(field string) {
	w.Str(field)
	w.byte(':')
}

// ObjEnd writes end of object token.
func (w *Writer) ObjEnd() {
	w.byte('}')
}

// ArrStart writes start of array.
func (w *Writer) ArrStart() {
	w.byte('[')
}

// ArrEnd writes end of array.
func (w *Writer) ArrEnd() {
	w.byte(']')
}

// Comma writes comma.
func (w *Writer) Comma() {
	w.byte(',')
}
