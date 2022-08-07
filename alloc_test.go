//go:build !race

package jx

import (
	"testing"

	"github.com/go-faster/errors"
)

const defaultAllocRuns = 20

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

func zeroAllocEnc(t *testing.T, f func(e *Encoder)) {
	t.Helper()
	var e Encoder
	zeroAlloc(t, func() {
		e.Reset()
		f(&e)
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
	t.Run("Encoder", func(t *testing.T) {
		t.Run("Manual", func(t *testing.T) {
			zeroAllocEnc(t, func(e *Encoder) {
				e.ObjStart()
				e.FieldStart("foo")
				e.ArrStart()
				e.Int(1)
				e.Int(2)
				e.Int(3)
				e.ArrEnd()
				e.ObjEnd()
			})
		})
		t.Run("Small object", func(t *testing.T) {
			zeroAllocEnc(t, encodeSmallObject)
		})
		t.Run("Callback", func(t *testing.T) {
			zeroAllocEnc(t, encodeSmallCallback)
		})
	})
}
