// Package jx implements encoding and decoding of json as per RFC 4627.
//
// The Iter provides a way to iterate over bytes/string/reader
// and yield parsed elements one by one, fast.
package jx

import "sync"

// Valid reports whether json in data is valid.
func Valid(data []byte) bool {
	i := GetIter()
	defer PutIter(i)
	i.ResetBytes(data)
	return i.Skip() == nil
}

var (
	streamPool = &sync.Pool{
		New: func() interface{} {
			return NewStream(nil, 256)
		},
	}
	iterPool = &sync.Pool{
		New: func() interface{} {
			return NewIter()
		},
	}
)

// GetIter gets *Iter from pool.
func GetIter() *Iter {
	return iterPool.Get().(*Iter)
}

// PutIter puts *Iter into pool.
func PutIter(i *Iter) {
	i.Reset(nil)
	iterPool.Put(i)
}

// GetStream returns *Stream from pool.
func GetStream() *Stream {
	return streamPool.Get().(*Stream)
}

// PutStream puts *Stream to pool
func PutStream(s *Stream) {
	s.Reset(nil)
	streamPool.Put(s)
}
