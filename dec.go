package jx

import (
	"io"

	"github.com/go-faster/errors"
)

// Type of json value.
type Type int

func (t Type) String() string {
	switch t {
	case Invalid:
		return "invalid"
	case String:
		return "string"
	case Number:
		return "number"
	case Null:
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
	// Invalid json value.
	Invalid Type = iota
	// String json value, like "foo".
	String
	// Number json value, like 100 or 1.01.
	Number
	// Null json value.
	Null
	// Bool json value, true or false.
	Bool
	// Array json value, like [1, 2, 3].
	Array
	// Object json value, like {"foo": 1}.
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
	for i := range types {
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
	types['n'] = Null
	types['['] = Array
	types['{'] = Object
}

// Decoder decodes json.
//
// Can decode from io.Reader or byte slice directly.
type Decoder struct {
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

const defaultBuf = 512

// Decode creates a Decoder that reads json from io.Reader.
func Decode(reader io.Reader, bufSize int) *Decoder {
	if bufSize <= 0 {
		bufSize = defaultBuf
	}
	return &Decoder{
		reader: reader,
		buf:    make([]byte, bufSize),
	}
}

// DecodeBytes creates a Decoder that reads json from byte slice.
func DecodeBytes(input []byte) *Decoder {
	return &Decoder{
		buf:  input,
		tail: len(input),
	}
}

// DecodeStr creates a Decoder that reads string as json.
func DecodeStr(input string) *Decoder {
	return DecodeBytes([]byte(input))
}

// Reset resets reader and underlying state, next reads will use provided io.Reader.
func (d *Decoder) Reset(reader io.Reader) {
	d.reader = reader
	d.head = 0
	d.tail = 0
	d.depth = 0

	// Reads from reader need buffer.
	if cap(d.buf) == 0 {
		// Allocate new buffer if none.
		d.buf = make([]byte, defaultBuf)
	}
	if len(d.buf) == 0 {
		// Set buffer to full capacity if needed.
		d.buf = d.buf[:cap(d.buf)]
	}
}

// ResetBytes resets underlying state, next reads will use provided buffer.
func (d *Decoder) ResetBytes(input []byte) {
	d.reader = nil
	d.head = 0
	d.tail = len(input)
	d.depth = 0

	d.buf = input
}

// Next gets Type of relatively next json element
func (d *Decoder) Next() Type {
	v, _ := d.next()
	d.unread()
	return types[v]
}

var spaceSet = [256]byte{
	' ': 1, '\n': 1, '\t': 1, '\r': 1,
}

func (d *Decoder) consume(c byte) (err error) {
	buf := d.buf[d.head:d.tail]
	var got byte
	if len(buf) > 0 && spaceSet[buf[0]] == 0 {
		d.head++
		got = buf[0]
	} else {
		got, err = d.more()
		if err != nil {
			return err
		}
	}
	if c != got {
		return badToken(got)
	}
	return nil
}

// more is next but io.EOF is unexpected.
func (d *Decoder) more() (byte, error) {
	c, err := d.next()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return c, err
}

// next reads next non-whitespace token or error.
func (d *Decoder) next() (byte, error) {
	for {
		buf := d.buf[d.head:d.tail]
		for i, c := range buf {
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			default:
				d.head += i + 1
				return c, nil
			}
		}
		if err := d.read(); err != nil {
			return 0, err
		}
	}
}

func (d *Decoder) byte() (byte, error) {
	if d.head == d.tail {
		err := d.read()
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			return 0, err
		}
	}
	c := d.buf[d.head]
	d.head++
	return c, nil
}

func (d *Decoder) read() error {
	if d.reader == nil {
		d.head = d.tail
		return io.EOF
	}

	n, err := d.reader.Read(d.buf)
	if err != nil {
		return err
	}

	d.head = 0
	d.tail = n
	return nil
}

func (d *Decoder) readAtLeast(min int) error {
	if d.reader == nil {
		d.head = d.tail
		return io.ErrUnexpectedEOF
	}

	if need := min - len(d.buf); need > 0 {
		d.buf = append(d.buf, make([]byte, need)...)
	}
	n, err := io.ReadAtLeast(d.reader, d.buf, min)
	if err != nil {
		if err == io.EOF && n == 0 {
			return io.ErrUnexpectedEOF
		}
		return err
	}

	d.head = 0
	d.tail = n
	return nil
}

func (d *Decoder) unread() { d.head-- }

// limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
const maxDepth = 10000

var errMaxDepth = errors.New("depth: maximum")

func (d *Decoder) incDepth() error {
	d.depth++
	if d.depth > maxDepth {
		return errMaxDepth
	}
	return nil
}

var errNegativeDepth = errors.New("depth: negative")

func (d *Decoder) decDepth() error {
	d.depth--
	if d.depth < 0 {
		return errNegativeDepth
	}
	return nil
}
