package jx

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_ObjBytes(t *testing.T) {
	t.Run("Object", func(t *testing.T) {
		i := DecodeStr(`{  "id" :1 ,  "randomNumber"  :  10    }`)
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
	})
	t.Run("Depth", func(t *testing.T) {
		var input []byte
		for i := 0; i <= maxDepth; i++ {
			input = append(input, `{"1":`...)
		}
		d := DecodeBytes(input)
		require.ErrorIs(t, d.ObjBytes(func(d *Decoder, key []byte) error {
			return crawlValue(d)
		}), errMaxDepth)
	})
	t.Run("Invalid", func(t *testing.T) {
		for _, s := range testObjs {
			checker := require.Error
			if json.Valid([]byte(s)) {
				continue
			}

			d := DecodeStr(s)
			err := d.ObjBytes(func(d *Decoder, key []byte) error {
				return crawlValue(d)
			})
			if err == nil && len(d.buf) > 0 {
				// FIXME(tdakkota): fix cases like {"hello":{}}}
				continue
			}
			checker(t, err, s)
		}
	})
}

func TestDecoderObjBytesIssue62(t *testing.T) {
	a := require.New(t)

	const input = `{"1":1,"2":2}`

	// Force decoder to read only first 4 bytes of input.
	d := Decode(strings.NewReader(input), 4)

	actual := map[string]int{}
	a.NoError(d.ObjBytes(func(d *Decoder, key []byte) error {
		val, err := d.Int()
		if err != nil {
			return err
		}
		actual[string(key)] = val
		return nil
	}))
	a.Equal(map[string]int{
		"1": 1,
		"2": 2,
	}, actual)
}
