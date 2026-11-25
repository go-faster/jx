package jx

import (
	"bytes"
	"io"
)

// Capture calls f and then rolls back to state before call.
func (d *Decoder) Capture(f func(d *Decoder) error) error {
	if f == nil {
		return nil
	}

	if d.reader != nil {
		var captured bytes.Buffer
		original := d.reader
		savedBuf := append([]byte(nil), d.buf...)
		streamOffset := d.streamOffset
		d.reader = io.TeeReader(original, &captured)
		defer func() {
			d.reader = io.MultiReader(&captured, original)
			d.buf = savedBuf
			d.streamOffset = streamOffset
		}()
	}
	head, tail, depth := d.head, d.tail, d.depth
	err := f(d)
	d.head, d.tail, d.depth = head, tail, depth
	return err
}
