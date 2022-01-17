package jx

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_SkipArrayNested(t *testing.T) {
	runTestCases(t, []string{
		`[-0.12, "stream"]`,
		`["hello", "stream"]`,
		`[null , "stream"]`,
		`[true , "stream"]`,
		`[false , "stream"]`,
		`[[1, [2, [3], 4]], "stream"]`,
		`[ [ ], "stream"]`,
	}, func(t *testing.T, d *Decoder) error {
		var err error
		a := require.New(t)
		_, err = d.Elem()
		a.NoError(err)
		err = d.Skip()
		a.NoError(err)
		_, err = d.Elem()
		a.NoError(err)
		if s, _ := d.Str(); s != "stream" {
			t.FailNow()
		}
		return nil
	})
}

func TestDecoderSkip_Nested(t *testing.T) {
	d := DecodeStr(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	if _, err := d.Elem(); err != nil {
		t.Fatal(err)
	}
	require.NoError(t, d.Skip())
	if _, err := d.Elem(); err != nil {
		t.Fatal(err)
	}
	s, err := d.Str()
	require.NoError(t, err)
	require.Equal(t, "stream", s)
}

func TestDecoderSkip_SimpleNested(t *testing.T) {
	d := DecodeStr(`["foo", "bar", "baz"]`)
	require.NoError(t, d.Skip())
}

func TestDecoder_skipNumber(t *testing.T) {
	inputs := []string{
		`0`,
		`120`,
		`0.`,
		`0.0e`,
		`0.0e+1`,
	}
	sr := strings.NewReader("")
	er := &errReader{}
	for i, tt := range inputs {
		t.Run(fmt.Sprintf("Test%d", i), func(t *testing.T) {
			sr.Reset(tt)
			d := Decode(io.MultiReader(sr, er), len(tt))
			require.NoError(t, d.read())
			require.Error(t, d.skipNumber())
		})
	}
}

func TestDecoder_SkipObjDepth(t *testing.T) {
	var input []byte
	for i := 0; i <= maxDepth; i++ {
		input = append(input, `{"1":`...)
	}
	require.Error(t, DecodeBytes(input).Skip())
}
