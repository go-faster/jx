// Package jx implements encoding and decoding of json as per RFC 4627.
//
// The Reader provides a way to iterate over bytes/string/reader
// and yield parsed elements one by one, fast.
package jx

import "sync"

// Valid reports whether json in data is valid.
func Valid(data []byte) bool {
	i := GetReader()
	defer PutReader(i)
	i.ResetBytes(data)
	return i.Skip() == nil
}

var (
	writePool = &sync.Pool{
		New: func() interface{} {
			return NewWriter(nil, 256)
		},
	}
	readPool = &sync.Pool{
		New: func() interface{} {
			return NewReader()
		},
	}
)

// GetReader gets *Reader from pool.
func GetReader() *Reader {
	return readPool.Get().(*Reader)
}

// PutReader puts *Reader into pool.
func PutReader(i *Reader) {
	i.Reset(nil)
	readPool.Put(i)
}

// GetWriter returns *Writer from pool.
func GetWriter() *Writer {
	return writePool.Get().(*Writer)
}

// PutWriter puts *Writer to pool
func PutWriter(s *Writer) {
	s.Reset(nil)
	writePool.Put(s)
}
