package jx

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

type limitWriter struct {
	w io.Writer
	n int64
}

func (t *limitWriter) Write(p []byte) (n int, err error) {
	if t.n-int64(len(p)) < 0 {
		return 0, errors.New("limit reached")
	}
	// real write
	n = len(p)
	if int64(n) > t.n {
		n = int(t.n)
	}
	n, err = t.w.Write(p[0:n])
	t.n -= int64(n)
	if err == nil {
		n = len(p)
	}
	return
}

func TestWriter_Base64(t *testing.T) {
	const bufSize = minEncoderBufSize
	const fieldLength = bufSize - len(`{"":`)

	limits := []int64{
		31, // write '"'
		32, // flush
		33, // Write base64
		73, // Write tail of base64
	}

	data := bytes.Repeat([]byte{0}, bufSize)
	for _, n := range limits {
		// Write '"' error.
		e := NewStreamingEncoder(&limitWriter{w: io.Discard, n: n}, bufSize)
		e.Obj(func(e *Encoder) {
			e.FieldStart(strings.Repeat("a", fieldLength))
			e.Base64(data)
		})
		require.Error(t, e.Close())
	}
}

func BenchmarkWriter_Base64(b *testing.B) {
	for _, n := range []int{
		128,
		256,
		512,
		1024,
	} {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			var v []byte
			for i := 0; i < n; i++ {
				v = append(v, byte(i%256))
			}

			b.ReportAllocs()
			b.SetBytes(int64(n))
			var w Writer
			for i := 0; i < b.N; i++ {
				w.Base64(v)
				w.Reset()
			}
		})
	}
}
