package jx

import (
	"golang.org/x/xerrors"
)

func (it *Iter) skipNumber() error {
	ok, err := it.skipNumberFast()
	if err != nil || ok {
		return err
	}
	it.unread()
	if _, err := it.Float64(); err != nil {
		return err
	}
	return nil
}

func (it *Iter) skipNumberFast() (ok bool, err error) {
	dotFound := false
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case tDot:
			if dotFound {
				return false, xerrors.New("more than one dot")
			}
			if i+1 == it.tail {
				return false, nil
			}
			c = it.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				return false, xerrors.New("no digit after dot")
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if it.head == i {
					return false, nil // if - without following digits
				}
				it.head = i
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

func (it *Iter) strSkip() error {
	ok, err := it.strFastSkip()
	if err != nil || ok {
		return err
	}

	it.unread()
	if _, err := it.str(value{ignore: true}); err != nil {
		return err
	}
	return nil
}

func (it *Iter) strFastSkip() (ok bool, err error) {
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		switch {
		case c == '"':
			it.head = i + 1
			return true, nil
		case c == '\\':
			return false, nil
		case c < ' ':
			return false, badToken(c)
		}
	}
	return false, nil
}

func (it *Iter) skipObject() error {
	it.unread()
	return it.ObjBytes(func(iter *Iter, _ []byte) error {
		return iter.Skip()
	})
}

func (it *Iter) skipArray() error {
	it.unread()
	return it.Array(func(iter *Iter) error {
		return iter.Skip()
	})
}
