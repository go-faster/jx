// Package jx implements encoding and decoding of json as per RFC 4627.
//
// The Decoder provides a way to iterate over bytes/string/reader
// and yield parsed elements one by one, fast.
package jx

import "sync"

// Valid reports whether json in data is valid.
func Valid(data []byte) bool {
	d := GetDecoder()
	defer PutDecoder(d)
	d.ResetBytes(data)
	return d.Skip() == nil
}

var (
	encPool = &sync.Pool{
		New: func() interface{} {
			return NewEncoder(nil, 256)
		},
	}
	decPool = &sync.Pool{
		New: func() interface{} {
			return NewDecoder()
		},
	}
)

// GetDecoder gets *Decoder from pool.
func GetDecoder() *Decoder {
	return decPool.Get().(*Decoder)
}

// PutDecoder puts *Decoder into pool.
func PutDecoder(i *Decoder) {
	i.Reset(nil)
	decPool.Put(i)
}

// GetWriter returns *Encoder from pool.
func GetWriter() *Encoder {
	return encPool.Get().(*Encoder)
}

// PutWriter puts *Encoder to pool
func PutWriter(s *Encoder) {
	s.Reset(nil)
	s.buf = s.buf[:0]
	encPool.Put(s)
}
