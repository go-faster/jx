package jx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestType_String(t *testing.T) {
	met := map[string]bool{}
	for i := Invalid; i <= Object+1; i++ {
		s := i.String()
		if s == "" {
			t.Error("blank")
		}
		if met[s] {
			t.Errorf("met %s", s)
		}
		met[s] = true
	}
	if len(met) != 8 {
		t.Error("unexpected met types")
	}
}

func TestDecoder_Reset(t *testing.T) {
	var d Decoder
	d.ResetBytes([]byte{})
	d.Reset(bytes.NewBufferString(`true`))
	v, err := d.Bool()
	require.NoError(t, err)
	require.True(t, v)
}
