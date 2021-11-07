package jx

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requireCompat fails if `encoding/json` will encode v differently than exp.
func requireCompat(t testing.TB, got []byte, v interface{}) {
	t.Helper()
	buf, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, string(buf), string(got))
}

func TestPutEncoder(t *testing.T) {
	var wg sync.WaitGroup
	for j := 0; j < 4; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1024; i++ {
				e := GetEncoder()
				e.Raw("false")
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
