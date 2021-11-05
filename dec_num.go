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
	var str bool
	switch d.Next() {
	case String:
		// Consume start of the string.
		d.head++
		str = true
	case Number: // float or integer
	default:
		return v, errors.Errorf("unexpected %s", d.Next())
	}
	if d.reader == nil {
		// Can use underlying buffer directly.
		v.Value = d.number()
	} else {
		buf, err := d.numberAppend(v.Value[:0])
		if err != nil {
			return v, errors.Wrap(err, "decode")
		}
		v.Value = buf
	}

	var dot bool
	for _, c := range v.Value {
		if c != '.' {
			continue
		}
		if dot {
			return v, errors.New("multiple dots in number")
		}
		dot = true
		break
	}
	if dot {
		v.Format = NumFormatFloat
		if str {
			v.Format = NumFormatFloatStr
		}
	} else {
		v.Format = NumFormatInt
		if str {
			v.Format = NumFormatIntStr
		}
	}

	// TODO(ernado): Additional validity checks

	return v, nil
}
