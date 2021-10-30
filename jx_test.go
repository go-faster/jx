package jx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireCompat fails if `encoding/json` will encode v differently than exp.
func requireCompat(t testing.TB, exp []byte, v interface{}) {
	t.Helper()
	buf, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, exp, buf)
}
