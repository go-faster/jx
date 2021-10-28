package test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	j "github.com/ogen-go/json"
)

func Test_invalid_float(t *testing.T) {
	inputs := []string{
		`1.e1`, // dot without following digit
		`1.`,   // dot can not be the last char
		``,     // empty number
		`01`,   // extra leading zero
		`-`,    // negative without digit
		`--`,   // double negative
		`--2`,  // double negative
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			should := require.New(t)
			iter := j.ParseString(j.ConfigDefault, input+",")
			iter.Skip()
			should.NotEqual(io.EOF, iter.Error)
			should.NotNil(iter.Error)
			iter = j.ParseString(j.ConfigDefault, input+",")
			iter.ReadFloat64()
			should.NotEqual(io.EOF, iter.Error)
			should.NotNil(iter.Error)
			iter = j.ParseString(j.ConfigDefault, input+",")
			iter.ReadFloat32()
			should.NotEqual(io.EOF, iter.Error)
			should.NotNil(iter.Error)
		})
	}
}

func Test_valid(t *testing.T) {
	should := require.New(t)
	should.True(j.Valid([]byte(`{}`)))
	should.False(j.Valid([]byte(`{`)))
}
