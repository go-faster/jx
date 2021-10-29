package jx

import (
	"io"
)

// IteratorPool is a thread safe pool of iterators with same configuration.
type IteratorPool interface {
	Iterator(data []byte) *Iterator
	PutIterator(iter *Iterator)
}

// StreamPool is a thread safe pool of streams with same configuration.
type StreamPool interface {
	Stream(writer io.Writer) *Stream
	PutStream(stream *Stream)
}

func (cfg *frozenConfig) Stream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) PutStream(stream *Stream) {
	stream.out = nil
	stream.Error = nil
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) Iterator(data []byte) *Iterator {
	iter := cfg.iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) PutIterator(iter *Iterator) {
	iter.ResetBytes(nil)
	iter.Reset(nil)
	cfg.iteratorPool.Put(iter)
}
