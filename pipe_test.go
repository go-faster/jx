package jx

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	// Test that Encoder and Decoder can communicate via pipe.
	r, w := io.Pipe()
	const objects = 1024 * 10
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() { _ = w.CloseWithError(io.EOF) }()
		e := NewEncoder()
		// Write objects to w.
		for i := 0; i < objects; i++ {
			e.Reset()
			e.ObjEmpty()
			if _, err := e.WriteTo(w); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	d := NewDecoder()
	d.Reset(r)
	// Read exact count of objects.
	for i := 0; i < objects; i++ {
		if err := d.Obj(nil); err != nil {
			t.Error(err)
			break
		}
	}

	// Assert correct end of pipe.
	assert.ErrorIs(t, d.Skip(), io.EOF, "unexpected read")
	assert.Equal(t, d.Next(), Invalid)

	// Wait for encoder to finish.
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}
