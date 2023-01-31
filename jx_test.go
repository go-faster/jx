package jx

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutEncoder(t *testing.T) {
	var wg sync.WaitGroup
	for j := 0; j < 4; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1024; i++ {
				e := GetEncoder()
				e.RawStr("false")
				assert.Equal(t, "false", e.String())
				PutEncoder(e)
			}
		}()
	}
	wg.Wait()
}

func TestPutDecoder(t *testing.T) {
	var wg sync.WaitGroup
	for j := 0; j < 4; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1024; i++ {
				d := GetDecoder()
				assert.Equal(t, d.Next(), Invalid)
				d.Reset(bytes.NewBufferString("false"))
				assert.Equal(t, d.Next(), Bool)
				v, err := d.Bool()
				assert.NoError(t, err)
				assert.Equal(t, d.Next(), Invalid)
				assert.False(t, v)
				PutDecoder(d)
			}
		}()
	}
	wg.Wait()
}
