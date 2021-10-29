package jx

import (
	"golang.org/x/xerrors"
)

func (it *Iterator) skipNumber() error {
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

func (it *Iterator) skipNumberFast() (ok bool, err error) {
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

func (it *Iterator) strSkip() error {
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

func (it *Iterator) strFastSkip() (ok bool, err error) {
	for i := it.head; i < it.tail; i++ {
		c := it.buf[i]
		if c == '"' {
			it.head = i + 1
			return true, nil
		} else if c == '\\' {
			return false, nil
		} else if c < ' ' {
			return false, badToken(c)
		}
	}
	return false, nil
}

func (it *Iterator) skipObject() error {
	it.unread()
	return it.ObjectBytes(func(iter *Iterator, _ []byte) error {
		return iter.Skip()
	})
}

func (it *Iterator) skipArray() error {
	it.unread()
	return it.Array(func(iter *Iterator) error {
		return iter.Skip()
	})
}
