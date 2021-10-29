package jx

import (
	"io"
)

// Stream writes json to io.Writer.
//
// Error is not returned as return value, but stored as Error member on this stream instance.
type Stream struct {
	cfg       *frozenConfig
	out       io.Writer
	buf       []byte
	indention int

	Error error
}

// NewStream create new stream instance.
// cfg can be jsoniter.ConfigDefault.
// out can be nil if write to internal buffer.
// bufSize is the initial size for the internal buffer in bytes.
func NewStream(cfg API, out io.Writer, bufSize int) *Stream {
	return &Stream{
		cfg:       cfg.(*frozenConfig),
		out:       out,
		buf:       make([]byte, 0, bufSize),
		Error:     nil,
		indention: 0,
	}
}

// Pool returns a pool can provide more stream with same configuration
func (s *Stream) Pool() StreamPool {
	return s.cfg
}

// Reset reuse this stream instance by assign a new writer
func (s *Stream) Reset(out io.Writer) {
	s.out = out
	s.buf = s.buf[:0]
}

// Available returns how many bytes are unused in the buffer.
func (s *Stream) Available() int {
	return cap(s.buf) - len(s.buf)
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (s *Stream) Buffered() int {
	return len(s.buf)
}

// Buffer if writer is nil, use this method to take the result.
func (s *Stream) Buffer() []byte {
	return s.buf
}

// SetBuffer allows to append to the internal buffer directly.
func (s *Stream) SetBuffer(buf []byte) {
	s.buf = buf
}

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (s *Stream) Write(p []byte) (nn int, err error) {
	s.buf = append(s.buf, p...)
	if s.out != nil {
		nn, err = s.out.Write(s.buf)
		s.buf = s.buf[nn:]
		return
	}
	return len(p), nil
}

// WriteByte writes a single byte.
func (s *Stream) writeByte(c byte) {
	s.buf = append(s.buf, c)
}

func (s *Stream) writeTwoBytes(c1 byte, c2 byte) {
	s.buf = append(s.buf, c1, c2)
}

func (s *Stream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	s.buf = append(s.buf, c1, c2, c3)
}

func (s *Stream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	s.buf = append(s.buf, c1, c2, c3, c4)
}

func (s *Stream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	s.buf = append(s.buf, c1, c2, c3, c4, c5)
}

// Flush writes any buffered data to the underlying io.Writer.
func (s *Stream) Flush() error {
	if s.out == nil {
		return nil
	}
	if s.Error != nil {
		return s.Error
	}
	_, err := s.out.Write(s.buf)
	if err != nil {
		if s.Error == nil {
			s.Error = err
		}
		return err
	}
	s.buf = s.buf[:0]
	return nil
}

// WriteRaw write string out without quotes, just like []byte.
func (s *Stream) WriteRaw(v string) {
	s.buf = append(s.buf, v...)
}

// WriteNil write null to stream.
func (s *Stream) WriteNil() {
	s.writeFourBytes('n', 'u', 'l', 'l')
}

// WriteTrue write true to stream.
func (s *Stream) WriteTrue() {
	s.writeFourBytes('t', 'r', 'u', 'e')
}

// WriteFalse write false to stream.
func (s *Stream) WriteFalse() {
	s.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

// WriteBool write true or false into stream.
func (s *Stream) WriteBool(val bool) {
	if val {
		s.WriteTrue()
	} else {
		s.WriteFalse()
	}
}

// WriteObjectStart write { with possible indention
func (s *Stream) WriteObjectStart() {
	s.indention += s.cfg.indentionStep
	s.writeByte('{')
	s.writeIndention(0)
}

// WriteObjectField write "field": with possible indention
func (s *Stream) WriteObjectField(field string) {
	s.WriteString(field)
	if s.indention > 0 {
		s.writeTwoBytes(':', ' ')
	} else {
		s.writeByte(':')
	}
}

// WriteObjectEnd write } with possible indention
func (s *Stream) WriteObjectEnd() {
	s.writeIndention(s.cfg.indentionStep)
	s.indention -= s.cfg.indentionStep
	s.writeByte('}')
}

// WriteEmptyObject write {}
func (s *Stream) WriteEmptyObject() {
	s.writeByte('{')
	s.writeByte('}')
}

// WriteMore write , with possible indention
func (s *Stream) WriteMore() {
	s.writeByte(',')
	s.writeIndention(0)
}

// WriteArrayStart write [ with possible indention
func (s *Stream) WriteArrayStart() {
	s.indention += s.cfg.indentionStep
	s.writeByte('[')
	s.writeIndention(0)
}

// WriteEmptyArray write []
func (s *Stream) WriteEmptyArray() {
	s.writeTwoBytes('[', ']')
}

// WriteArrayEnd write ] with possible indention
func (s *Stream) WriteArrayEnd() {
	s.writeIndention(s.cfg.indentionStep)
	s.indention -= s.cfg.indentionStep
	s.writeByte(']')
}

func (s *Stream) writeIndention(delta int) {
	if s.indention == 0 {
		return
	}
	s.writeByte('\n')
	toWrite := s.indention - delta
	for i := 0; i < toWrite; i++ {
		s.buf = append(s.buf, ' ')
	}
}
