// Package jx implements encoding and decoding of json as per RFC 4627.
//
// The Decoder provides a way to iterate over bytes/string/reader
// and yield parsed elements one by one, fast.
package jx

import (
	"io"
	"sync"
)

// Valid reports whether data is valid json.
func Valid(data []byte) bool {
	d := GetDecoder()
	defer PutDecoder(d)
	d.ResetBytes(data)

	// First encountered value skip should consume all buffer.
	if err := d.Skip(); err != nil {
		return false
	}
	// Check for any trialing json.
	if err := d.Skip(); err != io.EOF {
		return false
	}

	return true
}

var (
	encPool = &sync.Pool{
		New: func() interface{} {
			return &Encoder{}
		},
	}
	decPool = &sync.Pool{
		New: func() interface{} {
			return &Decoder{}
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

// GetEncoder returns *Encoder from pool.
func GetEncoder() *Encoder {
	return encPool.Get().(*Encoder)
}

// PutEncoder puts *Encoder to pool
func PutEncoder(e *Encoder) {
	e.Reset()
	e.SetIdent(0)
	encPool.Put(e)
}
