package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteNull(t *testing.T) {
	testEncoderModes(t, func(e *Encoder) {
		e.Null()
	}, "null")
}

func TestDecodeNullArrayElement(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.Str()
	should.NoError(err)
	should.Equal("a", s)
}

func TestDecodeNullString(t *testing.T) {
	should := require.New(t)
	iter := DecodeStr(`[null,"a"]`)
	should.True(iter.Elem())
	should.NoError(iter.Null())
	should.True(iter.Elem())
	s, err := iter.Str()
	should.NoError(err)
	should.Equal("a", s)
}

func TestDecodeNullSkip(t *testing.T) {
	iter := DecodeStr(`[null,"a"]`)
	iter.Elem()
	iter.Skip()
	iter.Elem()
	if s, _ := iter.Str(); s != "a" {
		t.FailNow()
	}
}

func TestNullError(t *testing.T) {
	a := require.New(t)
	var (
		b     = [4]byte{'n', 'u', 'l', 'l'}
		valid = b
	)
	for i := range b {
		// Reset buffer.
		b = valid
		for c := byte(0); c < 255; c++ {
			// Skip expected value.
			if valid[i] == c {
				continue
			}
			// Skip space as first character.
			if i == 0 && spaceSet[c] == 1 {
				continue
			}
			b[i] = c
			var token *badTokenErr
			a.ErrorAs(DecodeBytes(b[:]).Null(), &token)
			a.Equalf(c, token.Token, "%c != %c (%q)", c, token.Token, b)
		}
	}
}
