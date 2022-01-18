package jx

import "github.com/go-faster/errors"

// ArrIter is decoding array iterator.
type ArrIter struct {
	d      *Decoder
	closed bool
	comma  bool
}

// ArrIter creates new array iterator.
func (d *Decoder) ArrIter() (ArrIter, error) {
	if err := d.consume('['); err != nil {
		return ArrIter{}, errors.Wrap(err, "start")
	}
	if err := d.incDepth(); err != nil {
		return ArrIter{}, errors.Wrap(err, "inc")
	}
	if _, err := d.more(); err != nil {
		return ArrIter{}, err
	}
	d.unread()
	return ArrIter{d: d}, nil
}

// Next consumes element and returns false, if there is no elements anymore.
func (i *ArrIter) Next() (bool, error) {
	if i.closed {
		return false, nil
	}

	dec := i.d
	c, err := dec.more()
	if err != nil {
		return false, err
	}
	if c == ']' {
		i.closed = true
		return false, dec.decDepth()
	}
	if i.comma {
		if c != ',' {
			return false, badToken(c)
		}
	} else {
		dec.unread()
	}
	i.comma = true
	return true, nil
}
