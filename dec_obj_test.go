package jx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_ObjectBytes(t *testing.T) {
	i := DecodeStr(`{"id":1,"randomNumber":10}`)
	met := map[string]struct{}{}
	require.NoError(t, i.ObjBytes(func(i *Decoder, key []byte) error {
		switch string(key) {
		case "id":
			v, err := i.Int64()
			assert.NoError(t, err)
			assert.Equal(t, int64(1), v)
			met["id"] = struct{}{}
		case "randomNumber":
			v, err := i.Int64()
			if err != nil {
				return err
			}
			assert.Equal(t, int64(10), v)
			met["randomNumber"] = struct{}{}
		}
		return nil
	}))
	if len(met) != 2 {
		t.Error("not all keys met")
	}
}
