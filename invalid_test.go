package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
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
			iter := ReadString(input + ",")
			should.Error(iter.Skip())
			iter = ReadString(input + ",")
			_, err := iter.Float64()
			should.Error(err)
			iter = ReadString(input + ",")
			_, err = iter.Float32()
			should.Error(err)
		})
	}
}

func Test_valid(t *testing.T) {
	should := require.New(t)
	should.True(Valid([]byte(`{}`)))
	should.False(Valid([]byte(`{`)))
}
