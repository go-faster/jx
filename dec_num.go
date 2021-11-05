package jx

import (
	"github.com/ogen-go/errors"
)

// Num decodes number.
func (d *Decoder) Num() (Num, error) {
	return d.NumTo(Num{})
}

// NumTo decodes number into Num.
func (d *Decoder) NumTo(v Num) (Num, error) {
	switch d.Next() {
	case String:
		// Consume start of the string.
		d.head++
	case Number: // float or integer
	default:
		return v, errors.Errorf("unexpected %s", d.Next())
	}
	if d.reader == nil {
		// Can use underlying buffer directly.
		v = d.number()
	} else {
		buf, err := d.numberAppend(v[:0])
		if err != nil {
			return v, errors.Wrap(err, "decode")
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

	return v, nil
}
