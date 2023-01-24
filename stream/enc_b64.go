package stream

import "github.com/segmentio/asm/base64"

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
	buf := e.w.buf
	switch {
	case encodedLen <= cap(buf):
		// Case 2: There is enough space in the buffer after flushing.
		if e.w.flush() {
			return true
		}
		fallthrough
	case len(buf)+encodedLen <= cap(buf):
		// Case 1: There is enough space in the buffer.
		base64.StdEncoding.Encode(buf[len(buf):len(buf)+encodedLen], data)
		e.w.buf = buf[:len(buf)+encodedLen]
	default:
		// Case 3: There is not enough space in the buffer.
		// We need to flush the buffer and then write the encoded data directly.
		if e.w.flush() {
			return true
		}

		// StdEncoding includes padding, so we can't just split the data into chunks.
		// FIXME(tdakkota): Is there a way to avoid this allocation?
		//  If we can't, we should at least use a pool.
		//  Or use streaming encoder from stdlib.
		//  Or remove this method.
		r := make([]byte, encodedLen)
		base64.StdEncoding.Encode(r, data)
		if e.w.writeBytes(r...) {
			return true
		}
	}

	return e.w.writeByte('"')
}
