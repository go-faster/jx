package jx

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testDecoderRaw(t *testing.T, raw func(d *Decoder) (Raw, error)) {
	tests := []struct {
		input     string
		typ       Type
		expectErr bool
	}{
		{`"foo"`, String, false},
		{`"foo\`, Invalid, true},

		{`10`, Number, false},
		{`1asf0`, Invalid, true},

		{`null`, Null, false},
		{`nul`, Invalid, true},

		{`true`, Bool, false},
		{`tru`, Invalid, true},

		{`[1, 2, 3, 4, 5]`, Array, false},
		{`[1, 2, 3, 4, 5}`, Invalid, true},

		{`{"foo":"bar"}`, Object, false},
		{`{"foo":"bar", "baz":"foobar"}`, Object, false},
		{`{"foo":"bar}`, Invalid, true},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), testBufferReader(tt.input, func(t *testing.T, d *Decoder) {
			a := require.New(t)
			raw, err := raw(d)
			if tt.expectErr {
				a.Error(err)
				return
			}
			a.NoError(err)
			a.Equal(tt.input, raw.String())
			a.Equal(tt.typ, raw.Type())
		}))
	}

	objectInput := func(input string) string {
		e := GetEncoder()
		defer PutEncoder(e)

		e.Obj(func(e *Encoder) {
			const length = 8
			for i := range [length]struct{}{} {
				e.FieldStart(fmt.Sprintf("skip%d", i))
				e.Str("it")
			}

			e.FieldStart("test")
			e.RawStr(input)

			for i := range [length]struct{}{} {
				e.FieldStart(fmt.Sprintf("skip%d", i+length))
				e.Str("it")
			}
		})

		return e.String()
	}
	t.Run("InsideObject", func(t *testing.T) {
		for i, tt := range tests {
			tt := tt
			input := objectInput(tt.input)
			t.Run(fmt.Sprintf("Test%d", i+1), testBufferReader(input, func(t *testing.T, d *Decoder) {
				a := require.New(t)

				err := d.ObjBytes(func(d *Decoder, key []byte) error {
					if string(key) != "test" {
						return d.Skip()
					}
					raw, err := raw(d)
					if err != nil {
						return err
					}
					a.Equal(tt.input, raw.String())
					a.Equal(tt.typ, raw.Type())
					return nil
				})

				if tt.expectErr {
					a.Error(err)
				} else {
					a.NoError(err)
				}
			}))
		}
	})

	t.Run("InsideCapture", func(t *testing.T) {
		for i, tt := range tests {
			tt := tt
			input := objectInput(tt.input)
			t.Run(fmt.Sprintf("Test%d", i+1), testBufferReader(input, func(t *testing.T, d *Decoder) {
				a := require.New(t)

				err := d.Capture(func(d *Decoder) error {
					return d.ObjBytes(func(d *Decoder, key []byte) error {
						if string(key) != "test" {
							return d.Skip()
						}
						raw, err := raw(d)
						if err != nil {
							return err
						}
						a.Equal(tt.input, raw.String())
						a.Equal(tt.typ, raw.Type())
						return nil
					})
				})

				if tt.expectErr {
					a.Error(err)
				} else {
					a.NoError(err)
				}
			}))
		}
	})

	t.Run("InsideArray", func(t *testing.T) {
		for i, tt := range tests {
			tt := tt
			input := fmt.Sprintf(`[%s]`, tt.input)
			t.Run(fmt.Sprintf("Test%d", i+1), testBufferReader(input, func(t *testing.T, d *Decoder) {
				a := require.New(t)

				err := d.Arr(func(d *Decoder) error {
					raw, err := raw(d)
					if err != nil {
						return err
					}
					a.Equal(tt.input, raw.String())
					a.Equal(tt.typ, raw.Type())
					return nil
				})

				if tt.expectErr {
					a.Error(err)
				} else {
					a.NoError(err)
				}
			}))
		}
	})
}

func TestDecoder_Raw(t *testing.T) {
	testDecoderRaw(t, (*Decoder).Raw)
}

func TestDecoder_RawAppend(t *testing.T) {
	testDecoderRaw(t, func(d *Decoder) (Raw, error) {
		return d.RawAppend(nil)
	})
}

func BenchmarkRaw_Type(b *testing.B) {
	v := Raw{'1'}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if v.Type() != Number {
			b.Fatal("invalid")
		}
	}
}

func BenchmarkDecoder_Raw(b *testing.B) {
	data := []byte(`{"foo": [1,2,3,4,5,6,7,8,9,10,11,12,13,14]}`)
	b.ReportAllocs()

	b.Run("Bytes", func(b *testing.B) {
		var d Decoder
		for i := 0; i < b.N; i++ {
			d.ResetBytes(data)
			raw, err := d.Raw()
			if err != nil {
				b.Fatal(err)
			}
			if len(raw) == 0 {
				b.Fatal("blank")
			}
		}
	})
	b.Run("Reader", func(b *testing.B) {
		var (
			d Decoder
			r = new(bytes.Reader)
		)
		for i := 0; i < b.N; i++ {
			r.Reset(data)
			d.Reset(r)

			raw, err := d.Raw()
			if err != nil {
				b.Fatal(err)
			}
			if len(raw) == 0 {
				b.Fatal("blank")
			}
		}
	})
}
