package jx

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuite(t *testing.T) {
	// https://github.com/nst/JSONTestSuite
	// By Nicolas Seriot (https://github.com/nst)
	dir := filepath.Join("testdata", "test_parsing")
	files, err := os.ReadDir(dir)
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

		name := filepath.Join(dir, f.Name())
		action := f.Name()[:2]

		t.Run(f.Name(), func(t *testing.T) {
			data, err := os.ReadFile(name)
			require.NoError(t, err)

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
