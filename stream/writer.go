package stream

import (
	"io"

	"github.com/go-faster/jx/internal/byteseq"
)

type writer[W io.Writer] struct {
	writer   W
	buf      []byte
	flushErr error
}

func (w *writer[W]) Reset(writer W) {
	w.writer = writer
	w.buf = w.buf[:0]
	w.flushErr = nil
}

func (w *writer[W]) flush() bool {
	if w.flushErr != nil {
		return true
	}

	n, err := w.writer.Write(w.buf)
	switch {
	case err != nil:
		w.flushErr = err
		return true
	case n != len(w.buf):
		w.flushErr = io.ErrShortWrite
		return true
	default:
		w.buf = w.buf[:0]
		return false
	}
}

func writeByteseq[S byteseq.Byteseq, W io.Writer](w *writer[W], s S) bool {
	for len(w.buf)+len(s) > cap(w.buf) {
		if w.flush() {
			return true
		}

		n := copy(w.buf[len(w.buf):cap(w.buf)], s)
		s = s[n:]
		w.buf = w.buf[:len(w.buf)+n]
	}
	w.buf = append(w.buf, s...)
	return false
}

func (w *writer[W]) writeString(s string) bool {
	return writeByteseq(w, s)
}

func (w *writer[W]) writeBytes(s ...byte) bool {
	return writeByteseq(w, s)
}

func (w *writer[W]) writeByte(b byte) bool {
	return w.writeBytes(b)
}
