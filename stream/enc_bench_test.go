package stream_test

import (
	"io"
	"math/rand"
	"testing"

	"github.com/go-faster/jx"
	"github.com/go-faster/jx/stream"
)

func encodeFloats(e *jx.Encoder, arr []float64) {
	e.ArrStart()
	for _, v := range arr {
		e.Float64(v)
	}
	e.ArrEnd()
}

func encodeFloatsStream[W io.Writer](e *stream.Encoder[W], arr []float64) bool {
	if e.ArrStart() {
		return true
	}
	for _, v := range arr {
		if e.Float64(v) {
			return true
		}
	}
	return e.ArrEnd()
}

func BenchmarkEncodeFloats(b *testing.B) {
	const N = 100_000
	arr := make([]float64, N)
	for i := 0; i < N; i++ {
		arr[i] = rand.NormFloat64()
	}
	size := func() int64 {
		var enc jx.Encoder
		encodeFloats(&enc, arr)
		return int64(len(enc.Bytes()))
	}()
	b.Logf("Size: %d bytes", size)

	b.Run("Buffered", func(b *testing.B) {
		b.SetBytes(size)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Notice: no buffer reuse.
				var enc jx.Encoder
				encodeFloats(&enc, arr)
			}
		})
	})
	b.Run("Stream", func(b *testing.B) {
		b.SetBytes(size)
		b.RunParallel(func(pb *testing.PB) {
			enc := stream.NewEncoder(io.Discard)
			for pb.Next() {
				enc.Reset(io.Discard)
				encodeFloatsStream(enc, arr)
			}
		})
	})
}
