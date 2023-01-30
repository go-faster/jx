package jx

import "io"

const (
	minEncoderBufSize = 32
	encoderBufSize    = 512
)

// NewStreamingEncoder creates new streaming encoder.
func NewStreamingEncoder(w io.Writer, bufSize int) *Encoder {
	if bufSize < minEncoderBufSize {
		bufSize = encoderBufSize
	}
	return &Encoder{
		w: Writer{
			Buf:    make([]byte, 0, bufSize),
			stream: newStreamState(w),
		},
	}
}

// Close flushes underlying buffer to writer in streaming mode.
// Otherwise, it does nothing.
func (e *Encoder) Close() error {
	return e.w.Close()
}
