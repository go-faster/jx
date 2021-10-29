package jx

import (
	"io"

	"golang.org/x/xerrors"
)

// Type the type for JSON element
type Type int

func (t Type) String() string {
	switch t {
	case Invalid:
		return "invalid"
	case String:
		return "string"
	case Number:
		return "number"
	case Nil:
		return "nil"
	case Bool:
		return "bool"
	case Array:
		return "array"
	case Object:
		return "object"
	default:
		return "unknown"
	}
}

const (
	// Invalid invalid JSON element
	Invalid Type = iota
	// String JSON element "string"
	String
	// Number JSON element 100 or 0.10
	Number
	// Nil JSON element null
	Nil
	// Bool JSON element true or false
	Bool
	// Array JSON element []
	Array
	// Object JSON element {}
	Object
)

var hexDigits []byte
var types []Type

func init() {
	hexDigits = make([]byte, 256)
	for i := 0; i < len(hexDigits); i++ {
		hexDigits[i] = 255
	}
	for i := '0'; i <= '9'; i++ {
		hexDigits[i] = byte(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		hexDigits[i] = byte((i - 'a') + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		hexDigits[i] = byte((i - 'A') + 10)
	}
	types = make([]Type, 256)
	for i := 0; i < len(types); i++ {
		types[i] = Invalid
	}
	types['"'] = String
	types['-'] = Number
	types['0'] = Number
	types['1'] = Number
	types['2'] = Number
	types['3'] = Number
	types['4'] = Number
	types['5'] = Number
	types['6'] = Number
	types['7'] = Number
	types['8'] = Number
	types['9'] = Number
	types['t'] = Bool
	types['f'] = Bool
	types['n'] = Nil
	types['['] = Array
	types['{'] = Object
}

// Iter is an io.Reader like object, with json specific read functions.
//
// Error is not returned as return value, but rather stored as Error field.
type Iter struct {
	cfg    *frozenConfig
	reader io.Reader

	// buf is current buffer.
	//
	// Contains full json if reader is nil or used as a read buffer
	// otherwise.
	buf  []byte
	head int // offset in buf to start of current json stream
	tail int // offset in buf to end of current json stream

	depth int
}

// NewIter creates an empty Iter instance
func NewIter(cfg API) *Iter {
	return &Iter{
		cfg: cfg.(*frozenConfig),
	}
}

// Parse creates an Iter instance from io.Reader
func Parse(cfg API, reader io.Reader, bufSize int) *Iter {
	return &Iter{
		cfg:    cfg.(*frozenConfig),
		reader: reader,
		buf:    make([]byte, bufSize),
	}
}

// ParseBytes creates an Iter instance from byte array
func ParseBytes(cfg API, input []byte) *Iter {
	return &Iter{
		cfg:    cfg.(*frozenConfig),
		reader: nil,
		buf:    input,
		head:   0,
		tail:   len(input),
		depth:  0,
	}
}

// ParseString creates an Iter instance from string
func ParseString(cfg API, input string) *Iter {
	return ParseBytes(cfg, []byte(input))
}

// Reset reuse iterator instance by specifying another reader
func (it *Iter) Reset(reader io.Reader) *Iter {
	it.reader = reader
	it.head = 0
	it.tail = 0
	it.depth = 0
	return it
}

// ResetBytes reuse iterator instance by specifying another byte array as input
func (it *Iter) ResetBytes(input []byte) *Iter {
	it.reader = nil
	it.buf = input
	it.head = 0
	it.tail = len(input)
	it.depth = 0
	return it
}

// Next gets Type of relatively next json element
func (it *Iter) Next() Type {
	v, _ := it.next()
	it.unread()
	return types[v]
}

func (it *Iter) expectNext(c byte) error {
	v, err := it.next()
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	if err != nil {
		return err
	}
	if c != v {
		return badToken(v)
	}
	return nil
}

// next returns non-whitespace token or error.
func (it *Iter) next() (byte, error) {
	for {
		for i := it.head; i < it.tail; i++ {
			c := it.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			it.head = i + 1
			return c, nil
		}
		if err := it.read(); err != nil {
			return 0, err
		}
	}
}

func (it *Iter) byte() (ret byte) {
	if it.head == it.tail {
		if it.read() == nil {
			ret = it.buf[it.head]
			it.head++
			return ret
		}
		return 0
	}
	ret = it.buf[it.head]
	it.head++
	return ret
}

func (it *Iter) read() error {
	if it.reader == nil {
		it.head = it.tail
		return io.EOF
	}

	n, err := it.reader.Read(it.buf)
	if err != nil {
		return err
	}

	it.head = 0
	it.tail = n
	return nil
}

func (it *Iter) unread() { it.head-- }

// limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
const maxDepth = 10000

func (it *Iter) incrementDepth() error {
	it.depth++
	if it.depth > maxDepth {
		return xerrors.New("max depth")
	}
	return nil
}

func (it *Iter) decrementDepth() error {
	it.depth--
	if it.depth < 0 {
		return xerrors.New("negative depth")
	}
	return nil
}
