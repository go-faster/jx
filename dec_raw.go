package jx

import "github.com/go-faster/errors"

// Raw is like Skip(), but saves and returns skipped value as raw json.
//
// Do not retain returned value, it references underlying buffer.
func (d *Decoder) Raw() (Raw, error) {
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
func (d *Decoder) RawAppend(buf Raw) (Raw, error) {
	raw, err := d.Raw()
	if err != nil {
		return nil, err
	}
	return append(buf, raw...), err
}

// Raw json value.
type Raw []byte

// Type of Raw json value.
func (r Raw) Type() Type {
	d := Decoder{buf: r, tail: len(r)}
	return d.Next()
}

func (r Raw) String() string { return string(r) }
