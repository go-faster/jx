package jx

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

func TestEncoderStreamingCheck(t *testing.T) {
	a := require.New(t)

	e := NewStreamingEncoder(io.Discard, -1)

	_, err := e.Write([]byte("hello"))
	a.ErrorIs(err, errStreaming)

	_, err = e.WriteTo(io.Discard)
	a.ErrorIs(err, errStreaming)

	a.PanicsWithError(errStreaming.Error(), func() {
		_ = e.String()
	})
}

type errWriter struct {
	err error
	n   int
}

func (e *errWriter) Write(p []byte) (int, error) {
	n := e.n
	if n <= 0 {
		n = len(p)
	}
	return n, e.err
}

func TestEncoder_Close(t *testing.T) {
	errTest := errors.New("test")

	t.Run("FlushErr", func(t *testing.T) {
		ew := &errWriter{err: errTest}
		e := NewStreamingEncoder(ew, -1)
		e.Null()

		require.ErrorIs(t, e.Close(), errTest)
	})
	t.Run("WriteErr", func(t *testing.T) {
		ew := &errWriter{err: errTest}
		e := NewStreamingEncoder(ew, minEncoderBufSize)
		e.Obj(func(e *Encoder) {
			e.FieldStart(strings.Repeat("a", minEncoderBufSize))
			e.Null()
		})

		require.ErrorIs(t, e.Close(), errTest)
	})
	t.Run("ShortWrite", func(t *testing.T) {
		ew := &errWriter{n: 1}
		e := NewStreamingEncoder(ew, -1)
		e.Null()

		require.ErrorIs(t, e.Close(), io.ErrShortWrite)
	})
	t.Run("OK", func(t *testing.T) {
		e := NewStreamingEncoder(io.Discard, -1)
		e.Null()

		require.NoError(t, e.Close())
	})
	t.Run("NoStreaming", func(t *testing.T) {
		var e Encoder
		e.Null()

		require.NoError(t, e.Close())
	})
}

func TestEncoder_ResetWriter(t *testing.T) {
	do := func(e *Encoder) {
		e.ObjStart()
		e.FieldStart(strings.Repeat("a", minEncoderBufSize))
		e.Null()
		e.ObjEnd()

		require.NoError(t, e.Close())
	}

	var e Encoder
	do(&e)
	expected := e.String()

	for range [3]struct{}{} {
		var got strings.Builder
		e.ResetWriter(&got)
		do(&e)
		require.Equal(t, expected, got.String())
	}
}

// This benchmark is used to measure the overhead of ignoring errors.
func BenchmarkSkipError(b *testing.B) {
	e := NewStreamingEncoder(io.Discard, minEncoderBufSize)
	e.w.stream.setError(errors.New("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		encodeObject(e)
	}
}

func TestStreamingEncoderBufferSize(t *testing.T) {
	e := NewStreamingEncoder(io.Discard, -1)
	assert.Equal(t, encoderBufSize, cap(e.w.Buf))
	e = NewStreamingEncoder(io.Discard, minEncoderBufSize-1)
	assert.Equal(t, minEncoderBufSize, cap(e.w.Buf))
}
