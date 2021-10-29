package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_byte_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(Default, nil, 1)
	stream.byte('1')
	should.Equal("1", string(stream.Buf()))
	should.Equal(1, len(stream.buf))
	stream.byte('2')
	should.Equal("12", string(stream.Buf()))
	should.Equal(2, len(stream.buf))
	stream.threeBytes('3', '4', '5')
	should.Equal("12345", string(stream.Buf()))
}

func TestStream_Raw_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(Default, nil, 1)
	stream.Raw("123")
	should.NoError(stream.Flush())
	should.Equal("123", string(stream.Buf()))
}

func TestStream_Str_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(Default, nil, 0)
	stream.Str("123")
	should.NoError(stream.Flush())
	should.Equal(`"123"`, string(stream.Buf()))
}

type NopWriter struct {
	bufferSize int
}

func (w *NopWriter) Write(p []byte) (n int, err error) {
	w.bufferSize = cap(p)
	return len(p), nil
}

func TestStream_Flush_should_stop_grow_buffer(t *testing.T) {
	// GetStream an array of a zillion zeros.
	writer := new(NopWriter)
	stream := NewStream(Default, writer, 512)
	stream.ArrStart()
	for i := 0; i < 10000000; i++ {
		stream.WriteInt(0)
		stream.More()
		_ = stream.Flush()
	}
	stream.WriteInt(0)
	stream.ArrEnd()

	// Confirm that the buffer didn't have to grow.
	should := require.New(t)

	// 512 is the internal buffer size set in NewEncoder
	//
	// Flush is called after each array element, so only the first 8 bytes of it
	// is ever used, and it is never extended. Capacity remains 512.
	should.Equal(512, writer.bufferSize)
}
