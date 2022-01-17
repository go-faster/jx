package jx

import "testing"

func TestDecoder_Bool(t *testing.T) {
	runTestCases(t, testBools, func(t *testing.T, d *Decoder) error {
		_, err := d.Bool()
		return err
	})
}
