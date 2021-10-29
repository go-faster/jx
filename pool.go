package jx

import (
	"io"
)

// IterPool is a thread safe pool of iterators with same configuration.
type IterPool interface {
	GetIter(data []byte) *Iter
	PutIter(iter *Iter)
}

// StreamPool is a thread safe pool of streams with same configuration.
type StreamPool interface {
	GetStream(writer io.Writer) *Stream
	PutStream(stream *Stream)
}

func (cfg *frozenConfig) GetStream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) PutStream(stream *Stream) {
	stream.out = nil
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) GetIter(data []byte) *Iter {
	iter := cfg.iteratorPool.Get().(*Iter)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) PutIter(iter *Iter) {
	iter.ResetBytes(nil)
	iter.Reset(nil)
	cfg.iteratorPool.Put(iter)
}
