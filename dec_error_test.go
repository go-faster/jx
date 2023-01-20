package jx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_badTokenErr_Error(t *testing.T) {
	e := &badTokenErr{
		Token:  'c',
		Offset: 10,
	}
	s := error(e).Error()
	require.Equal(t, "unexpected byte 99 'c' at 10", s)
}
