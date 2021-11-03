package jx

import "github.com/ogen-go/errors"

// Raw is like Skip(), but saves and returns skipped value as raw json.
//
// Do not retain returned slice, it references underlying buffer.
func (d *Decoder) Raw() ([]byte, error) {
	if d.reader != nil {
		return nil, errors.New("not implemented for io.Reader")
	}

	start := d.head
	if err := d.Skip(); err != nil {
		return nil, errors.Wrap(err, "skip")
	}

	return d.buf[start:d.head], nil
}

// RawAppend is Raw that appends saved raw json value to buf.
func (d *Decoder) RawAppend(buf []byte) ([]byte, error) {
	raw, err := d.Raw()
	if err != nil {
		return nil, err
	}
	return append(buf, raw...), err
}
