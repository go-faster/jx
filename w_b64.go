package jx

import "github.com/segmentio/asm/base64"

// Base64 encodes data as standard base64 encoded string.
//
// Same as encoding/json, base64.StdEncoding or RFC 4648.
func (w *Writer) Base64(data []byte) {
	if data == nil {
		w.Null()
		return
	}

	w.byte('"')
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	start := len(w.Buf)
	w.Buf = append(w.Buf, make([]byte, encodedLen)...)
	base64.StdEncoding.Encode(w.Buf[start:], data)
	w.byte('"')
}
