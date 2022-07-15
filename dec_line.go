package jx

import "bytes"

var newLine = []byte{'\n'}

// Line returns current line number starting from 1.
func (d *Decoder) Line() int {
	return d.line + bytes.Count(d.buf[:d.head], newLine) + 1
}
