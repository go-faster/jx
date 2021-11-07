package jx

import (
	hexEnc "encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-faster/errors"
)

// NB: Any left intentionally unexported

// AnyType is type of Any value.
type AnyType byte

// Possible types for Any.
const (
	AnyInvalid AnyType = iota
	AnyStr
	AnyNumber
	AnyNull
	AnyObj
	AnyArr
	AnyBool
)

// Any represents any json value as sum type.
type Any struct {
	Type AnyType // zero value if AnyInvalid, can be AnyNull

	Str    string // AnyStr
	Bool   bool   // AnyBool
	Number Num    // AnyNumber

	// Key in object. Valid only if KeyValid.
	Key string
	// KeyValid denotes whether Any is element of object.
	// Needed for representing Key that is blank.
	//
	// Can be true only for Child of AnyObj.
	KeyValid bool

	Child []Any // AnyArr or AnyObj
}

// Equal reports whether v is equal to b.
func (v Any) Equal(b Any) bool {
	if v.KeyValid && v.Key != b.Key {
		return false
	}
	if v.Type != b.Type {
		return false
	}
	switch v.Type {
	case AnyNull, AnyInvalid:
		return true
	case AnyBool:
		return v.Bool == b.Bool
	case AnyStr:
		return v.Str == b.Str
	case AnyNumber:
		return v.Number.Equal(b.Number)
	}
	if len(v.Child) != len(b.Child) {
		return false
	}
	for i := range v.Child {
		if !v.Child[i].Equal(b.Child[i]) {
			return false
		}
	}
	return true
}

// Any reads Any value.
func (d *Decoder) Any() (Any, error) {
	var v Any
	if err := v.Read(d); err != nil {
		return Any{}, err
	}
	return v, nil
}

// Any writes Any value.
func (e *Encoder) Any(a Any) {
	a.Write(e)
}

func (v *Any) Read(d *Decoder) error {
	switch d.Next() {
	case Invalid:
		return errors.New("invalid")
	case Number:
		n, err := d.Num()
		if err != nil {
			return errors.Wrap(err, "number")
		}
		v.Number = n
		v.Type = AnyNumber
	case String:
		s, err := d.Str()
		if err != nil {
			return errors.Wrap(err, "str")
		}
		v.Str = s
		v.Type = AnyStr
	case Nil:
		if err := d.Null(); err != nil {
			return errors.Wrap(err, "null")
		}
		v.Type = AnyNull
	case Bool:
		b, err := d.Bool()
		if err != nil {
			return errors.Wrap(err, "bool")
		}
		v.Bool = b
		v.Type = AnyBool
	case Object:
		v.Type = AnyObj
		if err := d.Obj(func(r *Decoder, s string) error {
			var elem Any
			if err := elem.Read(r); err != nil {
				return errors.Wrap(err, "elem")
			}
			elem.Key = s
			elem.KeyValid = true
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return errors.Wrap(err, "obj")
		}
		return nil
	case Array:
		v.Type = AnyArr
		if err := d.Arr(func(r *Decoder) error {
			var elem Any
			if err := elem.Read(r); err != nil {
				return errors.Wrap(err, "elem")
			}
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return errors.Wrap(err, "array")
		}
		return nil
	}
	return nil
}

// Write json representation of Any to Encoder.
func (v Any) Write(w *Encoder) {
	if v.KeyValid {
		w.Field(v.Key)
	}
	switch v.Type {
	case AnyStr:
		w.Str(v.Str)
	case AnyNumber:
		w.Num(v.Number)
	case AnyBool:
		w.Bool(v.Bool)
	case AnyNull:
		w.Null()
	case AnyArr:
		w.ArrStart()
		for _, c := range v.Child {
			c.Write(w)
		}
		w.ArrEnd()
	case AnyObj:
		w.ObjStart()
		for _, c := range v.Child {
			c.Write(w)
		}
		w.ObjEnd()
	}
}

func (v Any) String() string {
	var b strings.Builder
	if v.KeyValid {
		if v.Key == "" {
			b.WriteString("<blank>")
		}
		b.WriteString(v.Key)
		b.WriteString(": ")
	}
	switch v.Type {
	case AnyStr:
		b.WriteString(`'` + v.Str + `'`)
	case AnyNumber:
		b.WriteString(v.Number.String())
	case AnyBool:
		b.WriteString(strconv.FormatBool(v.Bool))
	case AnyNull:
		b.WriteString("null")
	case AnyArr:
		b.WriteString("[")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("]")
	case AnyObj:
		b.WriteString("{")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("}")
	default:
		b.WriteString("<invalid>")
	}
	return b.String()
}

// Reset Any value to reuse.
func (v *Any) Reset() {
	v.Type = AnyInvalid
	v.Child = v.Child[:0]
	v.KeyValid = false

	v.Str = ""
	v.Key = ""
}

// Obj calls f for any child that is field if v is AnyObj.
func (v Any) Obj(f func(k string, v Any)) {
	if v.Type != AnyObj {
		return
	}
	for _, c := range v.Child {
		if !c.KeyValid {
			continue
		}
		f(c.Key, c)
	}
}

func TestAny_Read(t *testing.T) {
	t.Run("Obj", func(t *testing.T) {
		var v Any
		const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null,"f":"s","t":true,"":""}}`
		r := DecodeStr(input)
		assert.NoError(t, v.Read(r))
		assert.Equal(t, "{foo: {bar: 1, baz: [1, 2, 3.14], 200: null, f: 's', t: true, <blank>: ''}}", v.String())

		e := GetEncoder()
		e.Any(v)
		require.Equal(t, input, e.String(), "encoded value should equal to input")
	})
	t.Run("Inputs", func(t *testing.T) {
		for _, tt := range []struct {
			Input string
		}{
			{Input: "1"},
			{Input: "0.0"},
		} {
			t.Run(tt.Input, func(t *testing.T) {
				input := []byte(tt.Input)
				r := DecodeBytes(input)
				v, err := r.Any()
				require.NoError(t, err)

				e := GetEncoder()
				v.Write(e)

				var otherValue Any
				r.ResetBytes(e.Bytes())

				if err := otherValue.Read(r); err != nil {
					t.Error(err)
					t.Log(hexEnc.Dump(input))
					t.Log(hexEnc.Dump(e.Bytes()))
				}

				require.True(t, otherValue.Equal(v))
			})
		}
	})
	t.Run("Negative", func(t *testing.T) {
		for _, s := range []string{
			`foo`,
			`bar`,
			`tier`,
			`{]`,
			`[}`,
			`[foo]`,
			`[tier]`,
			`{foo:"`,
			`{"foo": tier`,
			`"baz`,
			`nil`,
		} {
			t.Run(s, func(t *testing.T) {
				d := DecodeStr(s)
				v, err := d.Any()
				require.Error(t, err)
				require.Equal(t, AnyInvalid, v.Type)
				require.Equal(t, "<invalid>", v.String())
			})
		}
		t.Run("Reader", func(t *testing.T) {
			d := Decode(errReader{}, -1)
			// Manually set internal buffer.
			d.tail = 1
			d.buf = []byte{'1'}

			// Trigger reading of number start from buffer
			// and call to reader.
			v, err := d.Any()
			require.Error(t, err)
			require.Equal(t, AnyInvalid, v.Type)
		})
	})
}

func TestAny_Equal(t *testing.T) {
	t.Run("ZeroValues", func(t *testing.T) {
		for _, typ := range []AnyType{
			AnyInvalid,
			AnyStr,
			AnyNumber,
			AnyNull,
			AnyObj,
			AnyArr,
			AnyBool,
		} {
			a := Any{Type: typ}
			t.Run("Equal", func(t *testing.T) {
				b := Any{Type: typ}
				require.True(t, a.Equal(b))
				t.Run("Child", func(t *testing.T) {
					aArr := Any{
						Type:  AnyArr,
						Child: []Any{a},
					}
					bArr := Any{
						Type: AnyArr,
					}
					require.False(t, aArr.Equal(bArr))
					bArr.Child = []Any{b}
					require.True(t, aArr.Equal(bArr))
				})
			})
			t.Run("NotEqual", func(t *testing.T) {
				b := Any{Type: typ + 1}
				require.False(t, a.Equal(b))
				t.Run("Child", func(t *testing.T) {
					aArr := Any{
						Type:  AnyArr,
						Child: []Any{a},
					}
					bArr := Any{
						Type:  AnyArr,
						Child: []Any{b},
					}
					require.False(t, aArr.Equal(bArr))
				})
			})
			t.Run("Keys", func(t *testing.T) {
				a.KeyValid = true
				b := Any{Type: typ, KeyValid: true}
				require.True(t, a.Equal(b))
				b.Key = "1"
				require.False(t, a.Equal(b))
			})
		}
	})
}

func BenchmarkAny(b *testing.B) {
	data := []byte(`[true, null, false, 100, "false"]`)
	r := GetDecoder()

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	var v Any
	for i := 0; i < b.N; i++ {
		v.Reset()
		r.ResetBytes(data)
		if err := v.Read(r); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnyStd(b *testing.B) {
	data := []byte(`[true, null, false, 100, "false"]`)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	var v []interface{}
	for i := 0; i < b.N; i++ {
		v = v[:0]
		if err := json.Unmarshal(data, &v); err != nil {
			b.Fatal(err)
		}
	}
}
