package jx

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter_Reset(t *testing.T) {
	w := GetWriter()
	defer PutWriter(w)

	w.True()
	require.NotEmpty(t, w.Buf)
	w.Reset()
	require.Empty(t, w.Buf)
}

func TestWriter_String(t *testing.T) {
	w := GetWriter()
	defer PutWriter(w)

	w.True()
	require.Equal(t, "true", w.String())
}

func TestWriter_Grow(t *testing.T) {
	should := require.New(t)
	e := &Writer{}
	should.Equal(0, len(e.Buf))
	should.Equal(0, cap(e.Buf))
	e.Grow(1024)
	should.Equal(0, len(e.Buf))
	should.Equal(1024, cap(e.Buf))
	e.Grow(512)
	should.Equal(0, len(e.Buf))
	should.Equal(1024, cap(e.Buf))
	e.Grow(4096)
	should.Equal(0, len(e.Buf))
	should.Equal(4096, cap(e.Buf))
}

func TestWriter_Write(t *testing.T) {
	t.Run("NonStreaming", func(t *testing.T) {
		// Test writing to buffer in non-streaming mode.
		w := &Writer{}
		data := []byte("test data")

		n, err := w.Write(data)

		require.NoError(t, err)
		require.Equal(t, len(data), n)
		require.Equal(t, data, w.Buf)
	})

	t.Run("NonStreaming_Multiple", func(t *testing.T) {
		// Test multiple writes accumulate in buffer.
		w := &Writer{}

		n1, err1 := w.Write([]byte("hello"))
		require.NoError(t, err1)
		require.Equal(t, 5, n1)

		n2, err2 := w.Write([]byte(" world"))
		require.NoError(t, err2)
		require.Equal(t, 6, n2)

		require.Equal(t, "hello world", string(w.Buf))
	})

	t.Run("Streaming_EmptyBuffer", func(t *testing.T) {
		// Test streaming mode without buffered data.
		var buf bytes.Buffer
		w := &Writer{}
		w.ResetWriter(&buf)

		data := []byte("streaming data")
		n, err := w.Write(data)

		require.NoError(t, err)
		require.Equal(t, len(data), n)
		require.Equal(t, "streaming data", buf.String())
	})

	t.Run("Streaming_WithBufferedData", func(t *testing.T) {
		// Test streaming mode with buffered data that needs flushing.
		var buf bytes.Buffer
		w := &Writer{}
		w.ResetWriter(&buf)

		// Add some data to the buffer
		w.Buf = append(w.Buf, []byte("buffered ")...)
		require.Greater(t, len(w.Buf), 0)

		// Now write more data
		data := []byte("new data")
		n, err := w.Write(data)

		require.NoError(t, err)
		require.Equal(t, len(data), n)
		// Buffer should be flushed and new data written
		require.Equal(t, "buffered new data", buf.String())
	})

	t.Run("Streaming_FlushError", func(t *testing.T) {
		// Test streaming mode when Flush fails (lines 19-22)
		errWriter := &errorWriter{err: io.ErrShortBuffer}
		w := &Writer{}
		w.ResetWriter(errWriter)

		// Add some data to the buffer to trigger flush
		w.Buf = append(w.Buf, []byte("data to flush")...)

		n, err := w.Write([]byte("more data"))

		require.Error(t, err)
		require.Equal(t, 0, n)
		require.Equal(t, io.ErrShortBuffer, err)
	})

	t.Run("Streaming_UnderlyingWriteError", func(t *testing.T) {
		// Test streaming mode when underlying writer fails (line 24)
		errWriter := &errorWriter{err: io.ErrClosedPipe}
		w := &Writer{}
		w.ResetWriter(errWriter)

		// No buffered data, so write goes directly to underlying writer
		n, err := w.Write([]byte("test"))

		require.Error(t, err)
		require.Equal(t, 0, n)
		require.Equal(t, io.ErrClosedPipe, err)
	})
}

// errorWriter is a test helper that simulates write errors
type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	return len(p), nil
}

func TestWriter_WriteTo(t *testing.T) {
	t.Run("NonStreaming_Success", func(t *testing.T) {
		// Test WriteTo in non-streaming mode with data in buffer.
		w := &Writer{}
		data := []byte("test data for WriteTo")
		w.Buf = append(w.Buf, data...)

		var buf bytes.Buffer
		n, err := w.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, int64(len(data)), n)
		require.Equal(t, string(data), buf.String())
	})

	t.Run("NonStreaming_EmptyBuffer", func(t *testing.T) {
		// Test WriteTo with empty buffer.
		w := &Writer{}

		var buf bytes.Buffer
		n, err := w.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, int64(0), n)
		require.Equal(t, "", buf.String())
	})

	t.Run("NonStreaming_LargeBuffer", func(t *testing.T) {
		// Test WriteTo with larger data to ensure all bytes are written.
		w := &Writer{}
		data := make([]byte, 10000)
		for i := range data {
			data[i] = byte(i % 256)
		}
		w.Buf = append(w.Buf, data...)

		var buf bytes.Buffer
		n, err := w.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, int64(len(data)), n)
		require.Equal(t, data, buf.Bytes())
	})

	t.Run("NonStreaming_WriteError", func(t *testing.T) {
		// Test WriteTo when the underlying writer returns an error.
		w := &Writer{}
		w.Buf = append(w.Buf, []byte("some data")...)

		errWriter := &errorWriter{err: io.ErrShortWrite}
		n, err := w.WriteTo(errWriter)

		require.Error(t, err)
		require.Equal(t, io.ErrShortWrite, err)
		require.Equal(t, int64(0), n)
	})

	t.Run("Streaming_ReturnsError", func(t *testing.T) {
		// Test WriteTo in streaming mode - should return errStreaming.
		var buf bytes.Buffer
		w := &Writer{}
		w.ResetWriter(&buf)

		// Add some data to buffer
		w.Buf = append(w.Buf, []byte("data")...)

		var outBuf bytes.Buffer
		n, err := w.WriteTo(&outBuf)

		require.Error(t, err)
		require.Equal(t, errStreaming, err)
		require.Equal(t, int64(0), n)
		require.Empty(t, outBuf.String())
	})

	t.Run("NonStreaming_JSONData", func(t *testing.T) {
		// Test WriteTo with actual JSON data built using Writer methods.
		w := &Writer{}
		w.ObjStart()
		w.FieldStart("name")
		w.Str("test")
		w.Comma()
		w.FieldStart("value")
		w.Int(42)
		w.ObjEnd()

		var buf bytes.Buffer
		n, err := w.WriteTo(&buf)

		require.NoError(t, err)
		require.Greater(t, n, int64(0))
		require.Equal(t, `{"name":"test","value":42}`, buf.String())
	})
}
