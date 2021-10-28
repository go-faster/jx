package jir

import (
	"fmt"
	"io"
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

// Iterator is a io.Reader like object, with JSON specific read functions.
// Error is not returned as return value, but stored as Error member on this iterator instance.
type Iterator struct {
	cfg              *frozenConfig
	reader           io.Reader
	buf              []byte
	head             int
	tail             int
	depth            int
	captureStartedAt int
	captured         []byte
	Error            error
	Attachment       interface{} // open for customized decoder
}

// NewIterator creates an empty Iterator instance
func NewIterator(cfg API) *Iterator {
	return &Iterator{
		cfg:    cfg.(*frozenConfig),
		reader: nil,
		buf:    nil,
		head:   0,
		tail:   0,
		depth:  0,
	}
}

// Parse creates an Iterator instance from io.Reader
func Parse(cfg API, reader io.Reader, bufSize int) *Iterator {
	return &Iterator{
		cfg:    cfg.(*frozenConfig),
		reader: reader,
		buf:    make([]byte, bufSize),
		head:   0,
		tail:   0,
		depth:  0,
	}
}

// ParseBytes creates an Iterator instance from byte array
func ParseBytes(cfg API, input []byte) *Iterator {
	return &Iterator{
		cfg:    cfg.(*frozenConfig),
		reader: nil,
		buf:    input,
		head:   0,
		tail:   len(input),
		depth:  0,
	}
}

// ParseString creates an Iterator instance from string
func ParseString(cfg API, input string) *Iterator {
	return ParseBytes(cfg, []byte(input))
}

// Pool returns a pool can provide more iterator with same configuration
func (iter *Iterator) Pool() IteratorPool {
	return iter.cfg
}

// Reset reuse iterator instance by specifying another reader
func (iter *Iterator) Reset(reader io.Reader) *Iterator {
	iter.reader = reader
	iter.head = 0
	iter.tail = 0
	iter.depth = 0
	iter.Error = nil
	return iter
}

// ResetBytes reuse iterator instance by specifying another byte array as input
func (iter *Iterator) ResetBytes(input []byte) *Iterator {
	iter.reader = nil
	iter.buf = input
	iter.head = 0
	iter.tail = len(input)
	iter.depth = 0
	iter.Error = nil
	return iter
}

// WhatIsNext gets Type of relatively next json element
func (iter *Iterator) WhatIsNext() Type {
	valueType := types[iter.nextToken()]
	iter.unreadByte()
	return valueType
}

func (iter *Iterator) nextToken() byte {
	// a variation of skip whitespaces, returning the next non-whitespace token
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			iter.head = i + 1
			return c
		}
		if !iter.loadMore() {
			return 0
		}
	}
}

// ReportError record a error in iterator instance with current position.
func (iter *Iterator) ReportError(operation string, msg string) {
	if iter.Error != nil {
		if iter.Error != io.EOF {
			return
		}
	}
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	peekEnd := iter.head + 10
	if peekEnd > iter.tail {
		peekEnd = iter.tail
	}
	parsing := string(iter.buf[peekStart:peekEnd])
	contextStart := iter.head - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := iter.head + 50
	if contextEnd > iter.tail {
		contextEnd = iter.tail
	}
	context := string(iter.buf[contextStart:contextEnd])
	iter.Error = fmt.Errorf("%s: %s, error found in #%v byte of ...|%s|..., bigger context ...|%s|...",
		operation, msg, iter.head-peekStart, parsing, context)
}

// CurrentBuffer gets current buffer as string for debugging purpose
func (iter *Iterator) CurrentBuffer() string {
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	return fmt.Sprintf("parsing #%v byte, around ...|%s|..., whole buffer ...|%s|...", iter.head,
		string(iter.buf[peekStart:iter.head]), string(iter.buf[0:iter.tail]))
}

func (iter *Iterator) readByte() (ret byte) {
	if iter.head == iter.tail {
		if iter.loadMore() {
			ret = iter.buf[iter.head]
			iter.head++
			return ret
		}
		return 0
	}
	ret = iter.buf[iter.head]
	iter.head++
	return ret
}

func (iter *Iterator) loadMore() bool {
	if iter.reader == nil {
		if iter.Error == nil {
			iter.head = iter.tail
			iter.Error = io.EOF
		}
		return false
	}
	if iter.captured != nil {
		iter.captured = append(iter.captured,
			iter.buf[iter.captureStartedAt:iter.tail]...)
		iter.captureStartedAt = 0
	}
	for {
		n, err := iter.reader.Read(iter.buf)
		if n == 0 {
			if err != nil {
				if iter.Error == nil {
					iter.Error = err
				}
				return false
			}
		} else {
			iter.head = 0
			iter.tail = n
			return true
		}
	}
}

func (iter *Iterator) unreadByte() {
	if iter.Error != nil {
		return
	}
	iter.head--
}

// limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
const maxDepth = 10000

func (iter *Iterator) incrementDepth() (success bool) {
	iter.depth++
	if iter.depth <= maxDepth {
		return true
	}
	iter.ReportError("incrementDepth", "exceeded max depth")
	return false
}

func (iter *Iterator) decrementDepth() (success bool) {
	iter.depth--
	if iter.depth >= 0 {
		return true
	}
	iter.ReportError("decrementDepth", "unexpected negative nesting")
	return false
}
