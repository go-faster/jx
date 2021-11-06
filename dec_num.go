package jx

import (
	"github.com/go-faster/errors"
)

// Num decodes number.
//
// Do not retain returned value, it references underlying buffer.
func (d *Decoder) Num() (Num, error) {
	return d.num(nil, false)
}

// NumAppend appends number.
func (d *Decoder) NumAppend(v Num) (Num, error) {
	return d.num(v, true)
}

// num decodes number.
func (d *Decoder) num(v Num, forceAppend bool) (Num, error) {
	var str bool
	switch d.Next() {
	case String:
		str = true
	case Number: // float or integer
	default:
		return v, errors.Errorf("unexpected %s", d.Next())
	}
	if d.reader == nil && !forceAppend {
		// Can use underlying buffer directly.
		start := d.head
		d.head++
		d.number()
		if str {
			if err := d.consume('"'); err != nil {
				return nil, errors.Wrap(err, "end of string")
			}
		}
		v = d.buf[start:d.head]
	} else {
		if str {
			d.head++ // '"'
			v = append(v, '"')
		}
		buf, err := d.numberAppend(v)
		if err != nil {
			return v, errors.Wrap(err, "decode")
		}
		if str {
			if err := d.consume('"'); err != nil {
				return nil, errors.Wrap(err, "end of string")
			}
			buf = append(buf, '"')
		}
		v = buf
	}

	var dot bool
	for _, c := range v {
		if c != '.' {
			continue
		}
		if dot {
			return v, errors.New("multiple dots in number")
		}
		dot = true
	}

	// TODO(ernado): Additional validity checks
	// Current invariants:
	// 1) Zero or one dot
	// 2) Only: +, -, ., e, E, 0-9

	return v, nil
}
