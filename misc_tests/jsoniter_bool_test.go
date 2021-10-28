package misc_tests

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ogen-go/json"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := json.ParseString(json.ConfigDefault, `true`)
	should.True(iter.ReadBool())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := json.ParseString(json.ConfigDefault, `false`)
	should.False(iter.ReadBool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := json.NewStream(json.ConfigDefault, buf, 4096)
	stream.WriteTrue()
	stream.WriteFalse()
	stream.WriteBool(false)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("truefalsefalse", buf.String())
}
