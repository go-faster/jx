package jx

import (
	"golang.org/x/xerrors"
)

func (r *Reader) skipNumber() error {
	ok, err := r.skipNumberFast()
	if err != nil || ok {
		return err
	}
	r.unread()
	if _, err := r.Float64(); err != nil {
		return err
	}
	return nil
}

func (r *Reader) skipNumberFast() (ok bool, err error) {
	dotFound := false
	for i := r.head; i < r.tail; i++ {
		c := r.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case tDot:
			if dotFound {
				return false, xerrors.New("more than one dot")
			}
			if i+1 == r.tail {
				return false, nil
			}
			c = r.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				return false, xerrors.New("no digit after dot")
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if r.head == i {
					return false, nil // if - without following digits
				}
				r.head = i
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

func (r *Reader) strSkip() error {
	ok, err := r.strFastSkip()
	if err != nil || ok {
		return err
	}

	r.unread()
	if _, err := r.str(value{ignore: true}); err != nil {
		return err
	}
	return nil
}

func (r *Reader) strFastSkip() (ok bool, err error) {
	for i := r.head; i < r.tail; i++ {
		c := r.buf[i]
		switch {
		case c == '"':
			r.head = i + 1
			return true, nil
		case c == '\\':
			return false, nil
		case c < ' ':
			return false, badToken(c)
		}
	}
	return false, nil
}

func (r *Reader) skipObject() error {
	r.unread()
	return r.ObjBytes(func(iter *Reader, _ []byte) error {
		return iter.Skip()
	})
}

func (r *Reader) skipArray() error {
	r.unread()
	return r.Array(func(iter *Reader) error {
		return iter.Skip()
	})
}
