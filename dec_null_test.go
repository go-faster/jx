package jx

import "testing"

func TestDecoder_Null(t *testing.T) {
	runTestCases(t, []string{
		"",
		"nope",
		"nul",
		"nil",
		"nul\x00",
		"null",
	}, func(t *testing.T, d *Decoder) error {
		return d.Null()
	})
}
