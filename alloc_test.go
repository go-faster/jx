//go:build !race && !gofuzz && go1.17
// +build !race,!gofuzz,go1.17

package jx

import (
	"testing"

	"github.com/go-faster/errors"
)

const defaultAllocRuns = 10

func zeroAlloc(t *testing.T, f func()) {
	t.Helper()
	avg := testing.AllocsPerRun(defaultAllocRuns, f)
	if avg > 0 {
		t.Errorf("Allocated %f bytes per run", avg)
	}
}

func zeroAllocDec(t *testing.T, buf []byte, f func(d *Decoder) error) {
	t.Helper()
	d := DecodeBytes(buf)
	zeroAlloc(t, func() {
		d.ResetBytes(buf)
		if err := f(d); err != nil {
			t.Fatal(err)
		}
	})
}

func zeroAllocDecStr(t *testing.T, s string, f func(d *Decoder) error) {
	t.Helper()
	zeroAllocDec(t, []byte(s), f)
}

func TestZeroAlloc(t *testing.T) {
	// Tests that checks for zero allocations.
	t.Run("Decoder", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			zeroAllocDec(t, benchData, func(d *Decoder) error {
				return d.Validate()
			})
		})
		t.Run("ObjBytes", func(t *testing.T) {
			zeroAllocDec(t, benchData, func(d *Decoder) error {
				return d.Arr(func(d *Decoder) error {
					return d.ObjBytes(func(d *Decoder, key []byte) error {
						switch string(key) {
						case "person", "company": // ok
						default:
							return errors.New("unexpected key")
						}
						switch d.Next() {
						case Object:
							return d.ObjBytes(func(d *Decoder, key []byte) error {
								return d.Skip()
							})
						default:
							return d.Skip()
						}
					})
				})
			})
		})
		t.Run("Int", func(t *testing.T) {
			zeroAllocDecStr(t, "12345", func(d *Decoder) error {
				v, err := d.Int()
				if v != 12345 {
					t.Fatal(v)
				}
				return err
			})
		})
		t.Run("StrBytes", func(t *testing.T) {
			zeroAllocDecStr(t, `"hello"`, func(d *Decoder) error {
				v, err := d.StrBytes()
				if string(v) != "hello" {
					t.Fatal(string(v))
				}
				return err
			})
		})
		t.Run("ArrBigFile", func(t *testing.T) {
			zeroAllocDec(t, benchData, func(d *Decoder) error {
				return d.Arr(nil)
			})
		})
	})
}
