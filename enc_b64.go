package jx

import "encoding/base64"

// Base64 encodes data as standard base64 encoded string.
//
// Same as encoding/json, base64.StdEncoding or RFC 4648.
func (e *Encoder) Base64(data []byte) {
	if data == nil {
		e.Null()
		return
	}
	e.comma()
	e.byte('"')
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	start := len(e.buf)
	e.buf = append(e.buf, make([]byte, encodedLen)...)
	base64.StdEncoding.Encode(e.buf[start:], data)
	e.byte('"')
}
