package stream

import (
	stdbase64 "encoding/base64"

	"github.com/segmentio/asm/base64"
)

// Base64 encodes data as standard base64 encoded string.
//
// Same as encoding/json, base64.StdEncoding or RFC 4648.
func (e *Encoder[W]) Base64(data []byte) bool {
	if data == nil {
		return e.Null()
	}

	if e.comma() || e.w.writeByte('"') {
		return true
	}

	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	switch {
	case encodedLen <= cap(e.w.buf):
		// Case 2: There is enough space in the buffer after flushing.
		if e.w.flush() {
			return true
		}
		fallthrough
	case len(e.w.buf)+encodedLen <= cap(e.w.buf):
		// Case 1: There is enough space in the buffer.
		base64.StdEncoding.Encode(e.w.buf[len(e.w.buf):len(e.w.buf)+encodedLen], data)
		e.w.buf = e.w.buf[:len(e.w.buf)+encodedLen]
	default:
		// Case 3: There is not enough space in the buffer.
		// We need to flush the buffer and then write the encoded data directly.
		if e.base64Slow(data) {
			return true
		}
	}

	return e.w.writeByte('"')
}

func (e *Encoder[W]) base64Slow(data []byte) bool {
	if e.w.flush() {
		return true
	}

	// StdEncoding includes padding, so we can't just split the data into chunks.
	// FIXME(tdakkota): Is there a way to avoid this?
	//  Can we use asm/base64 for this?
	//  Or remove this method.
	enc := stdbase64.NewEncoder(stdbase64.StdEncoding, e.w.writer)
	if _, err := enc.Write(data); err != nil {
		return e.w.checkErr(err)
	}
	return e.w.checkErr(enc.Close())
}
