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

// Reader is an io.Reader like object, with json specific read functions.
//
// Error is not returned as return value, but rather stored as Error field.
type Reader struct {
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

// NewReader creates an empty Reader instance.
func NewReader() *Reader {
	return &Reader{}
}

// Read creates an Reader instance from io.Reader
func Read(reader io.Reader, bufSize int) *Reader {
	return &Reader{
		reader: reader,
		buf:    make([]byte, bufSize),
	}
}

// ReadBytes creates a Reader instance from byte slice.
func ReadBytes(input []byte) *Reader {
	return &Reader{
		reader: nil,
		buf:    input,
		head:   0,
		tail:   len(input),
		depth:  0,
	}
}

// ReadString creates a Reader instance from string.
func ReadString(input string) *Reader {
	return ReadBytes([]byte(input))
}

// Reset reuse iterator instance by specifying another reader
func (r *Reader) Reset(reader io.Reader) *Reader {
	r.reader = reader
	r.head = 0
	r.tail = 0
	r.depth = 0
	return r
}

// ResetBytes reuse iterator instance by specifying another byte array as input
func (r *Reader) ResetBytes(input []byte) *Reader {
	r.reader = nil
	r.buf = input
	r.head = 0
	r.tail = len(input)
	r.depth = 0
	return r
}

// Next gets Type of relatively next json element
func (r *Reader) Next() Type {
	v, _ := r.next()
	r.unread()
	return types[v]
}

func (r *Reader) expectNext(c byte) error {
	v, err := r.next()
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
func (r *Reader) next() (byte, error) {
	for {
		for i := r.head; i < r.tail; i++ {
			c := r.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			r.head = i + 1
			return c, nil
		}
		if err := r.read(); err != nil {
			return 0, err
		}
	}
}

func (r *Reader) byte() (byte, error) {
	if r.head == r.tail {
		if err := r.read(); err != nil {
			return 0, err
		}
	}
	c := r.buf[r.head]
	r.head++
	return c, nil
}

func (r *Reader) read() error {
	if r.reader == nil {
		r.head = r.tail
		return io.EOF
	}

	n, err := r.reader.Read(r.buf)
	if err != nil {
		return err
	}

	r.head = 0
	r.tail = n
	return nil
}

func (r *Reader) unread() { r.head-- }

// limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
const maxDepth = 10000

func (r *Reader) incrementDepth() error {
	r.depth++
	if r.depth > maxDepth {
		return xerrors.New("max depth")
	}
	return nil
}

func (r *Reader) decrementDepth() error {
	r.depth--
	if r.depth < 0 {
		return xerrors.New("negative depth")
	}
	return nil
}
