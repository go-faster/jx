package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_byte_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.byte('1')
	should.Equal("1", string(e.Bytes()))
	should.Equal(1, len(e.buf))
	e.byte('2')
	should.Equal("12", string(e.Bytes()))
	should.Equal(2, len(e.buf))
	e.threeBytes('3', '4', '5')
	should.Equal("12345", string(e.Bytes()))
}

func TestEncoder(t *testing.T) {
	data := `"hello world"`
	buf := []byte(data)
	var e Encoder
	t.Run("Write", func(t *testing.T) {
		e.Reset()
		n, err := e.Write(buf)
		require.NoError(t, err)
		require.Equal(t, n, len(buf))
		require.Equal(t, data, e.String())
	})
	t.Run("RawBytes", func(t *testing.T) {
		e.Reset()
		e.RawBytes(buf)
		require.Equal(t, data, e.String())
	})
	t.Run("SetBytes", func(t *testing.T) {
		e.Reset()
		e.SetBytes(buf)
		require.Equal(t, data, e.String())
	})
}

func TestEncoder_Raw_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.Raw("123")
	should.Equal("123", string(e.Bytes()))
}

func TestEncoder_Str_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	e := GetEncoder()
	e.Str("123")
	should.Equal(`"123"`, string(e.Bytes()))
}

func TestEncoder_ArrEmpty(t *testing.T) {
	e := GetEncoder()
	e.ArrEmpty()
	require.Equal(t, "[]", string(e.Bytes()))
}

func TestEncoder_ObjEmpty(t *testing.T) {
	e := GetEncoder()
	e.ObjEmpty()
	require.Equal(t, "{}", string(e.Bytes()))
}
