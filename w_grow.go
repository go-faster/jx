package jx 

import "bytes"

// Grow grows the underlying buffer.
// It calls (*bytes.Buffer).Grow(n int) on b.Buf
func (b *Buffer) Grow(n int) {
	buf := bytes.NewBuffer(b.Buf)
	buf.Grow(n)
	b.Buf = buf.Bytes()
}
