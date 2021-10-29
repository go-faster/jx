package jx

import (
	"io"
)

// Stream writes json to io.Writer.
//
// Error is not returned as return value, but stored as Error member on this stream instance.
type Stream struct {
	out io.Writer
	buf []byte

	ident    int
	curIdent int
}

// SetIdent sets length of single indentation step.
func (s *Stream) SetIdent(n int) {
	s.ident = n
}

// NewStream create new stream instance.
// cfg can be jsoniter.ConfigDefault.
// out can be nil if write to internal buffer.
// bufSize is the initial size for the internal buffer in bytes.
func NewStream(out io.Writer, bufSize int) *Stream {
	return &Stream{
		out: out,
		buf: make([]byte, 0, bufSize),
	}
}

// Reset reuse this stream instance by assign a new writer
func (s *Stream) Reset(out io.Writer) {
	s.out = out
	s.buf = s.buf[:0]
}

// Buf returns underlying buffer.
func (s *Stream) Buf() []byte { return s.buf }

// SetBuf allows to set the internal buffer directly.
func (s *Stream) SetBuf(buf []byte) { s.buf = buf }

// byte writes a single byte.
func (s *Stream) byte(c byte) {
	s.buf = append(s.buf, c)
}

func (s *Stream) twoBytes(c1, c2 byte) {
	s.buf = append(s.buf, c1, c2)
}

func (s *Stream) threeBytes(c1, c2, c3 byte) {
	s.buf = append(s.buf, c1, c2, c3)
}

func (s *Stream) fourBytes(c1, c2, c3, c4 byte) {
	s.buf = append(s.buf, c1, c2, c3, c4)
}

func (s *Stream) fiveBytes(c1, c2, c3, c4, c5 byte) {
	s.buf = append(s.buf, c1, c2, c3, c4, c5)
}

// Flush writes any buffered data to the underlying io.Writer.
func (s *Stream) Flush() error {
	if s.out == nil {
		return nil
	}
	if _, err := s.out.Write(s.buf); err != nil {
		return err
	}
	s.buf = s.buf[:0]
	return nil
}

// Raw write string out without quotes, just like []byte.
func (s *Stream) Raw(v string) {
	s.buf = append(s.buf, v...)
}

// Null write null to stream.
func (s *Stream) Null() {
	s.fourBytes('n', 'u', 'l', 'l')
}

// True write true to stream.
func (s *Stream) True() {
	s.fourBytes('t', 'r', 'u', 'e')
}

// False writes false to stream.
func (s *Stream) False() {
	s.fiveBytes('f', 'a', 'l', 's', 'e')
}

// Bool writes boolean.
func (s *Stream) Bool(val bool) {
	if val {
		s.True()
	} else {
		s.False()
	}
}

// ObjStart writes { with possible indention.
func (s *Stream) ObjStart() {
	s.curIdent += s.ident
	s.byte('{')
	s.writeIdent(0)
}

// ObjField write "field": with possible indention.
func (s *Stream) ObjField(field string) {
	s.Str(field)
	if s.curIdent > 0 {
		s.twoBytes(':', ' ')
	} else {
		s.byte(':')
	}
}

// ObjEnd write } with possible indention
func (s *Stream) ObjEnd() {
	s.writeIdent(s.ident)
	s.curIdent -= s.ident
	s.byte('}')
}

// ObjEmpty write {}
func (s *Stream) ObjEmpty() {
	s.byte('{')
	s.byte('}')
}

// More write , with possible indention
func (s *Stream) More() {
	s.byte(',')
	s.writeIdent(0)
}

// ArrStart writes [ with possible indention.
func (s *Stream) ArrStart() {
	s.curIdent += s.ident
	s.byte('[')
	s.writeIdent(0)
}

// ArrEmpty writes [].
func (s *Stream) ArrEmpty() {
	s.twoBytes('[', ']')
}

// ArrEnd writes ] with possible indention.
func (s *Stream) ArrEnd() {
	s.writeIdent(s.ident)
	s.curIdent -= s.ident
	s.byte(']')
}

func (s *Stream) writeIdent(delta int) {
	if s.curIdent == 0 {
		return
	}
	s.byte('\n')
	toWrite := s.curIdent - delta
	for i := 0; i < toWrite; i++ {
		s.buf = append(s.buf, ' ')
	}
}
