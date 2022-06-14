package jx

import (
	"bytes"
	"io"

	"github.com/go-faster/errors"
)

// Raw is like Skip(), but saves and returns skipped value as raw json.
//
// Do not retain returned value, it references underlying buffer.
func (d *Decoder) Raw() (Raw, error) {
	start := d.head
	if orig := d.reader; orig != nil {
		buf := bytes.Buffer{}
		buf.Write(d.buf[d.head:d.tail])
		d.reader = io.TeeReader(orig, &buf)
		defer func() {
			d.reader = orig
		}()

		if err := d.Skip(); err != nil {
			return nil, errors.Wrap(err, "skip")
		}

		unread := d.tail - d.head
		raw := buf.Bytes()
		raw = raw[:len(raw)-unread]
		return raw, nil
	}

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
