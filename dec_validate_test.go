package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_Validate(t *testing.T) {
	d := GetDecoder()
	runTestdata(t.Fatal, func(name string, data []byte) {
		t.Run(name, func(t *testing.T) {
			d.ResetBytes(data)
			require.NoError(t, d.Validate())
		})
	})
}
