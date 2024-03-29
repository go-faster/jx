package jx

import (
	"encoding/json"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuite(t *testing.T) {
	// https://github.com/nst/JSONTestSuite
	// By Nicolas Seriot (https://github.com/nst)
	dir := path.Join("testdata", "test_parsing")
	files, err := testdata.ReadDir(dir)
	require.NoError(t, err)

	const (
		Accept    = "y_"
		Reject    = "n_"
		Undefined = "i_"
	)
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(f.Name(), ".json")
		action := f.Name()[:2]

		file := path.Join(dir, f.Name())
		data, err := testdata.ReadFile(file)
		require.NoError(t, err)

		t.Run(name, func(t *testing.T) {
			r := GetDecoder()
			r.ResetBytes(data)
			defer PutDecoder(r)

			_, decodeErr := r.Any()
			switch action {
			case Accept:
				assert.True(t, Valid(data), "validate")
				assert.True(t, json.Valid(data), "std")
				assert.NoError(t, decodeErr, "%#v", string(data))
			case Reject:
				assert.False(t, Valid(data), "validate")
				assert.False(t, json.Valid(data), "std")
				// TODO: assert decodeErr + buffer drain?
			case Undefined:
				if decodeErr == nil {
					t.Log("Accept")
				} else {
					t.Logf("Reject: %v", decodeErr)
				}
			default:
				t.Fatalf("Unknown prefix %q", action)
			}
		})
	}
}
