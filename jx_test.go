package jx

import (
	"bytes"
	hexEnc "encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func TestJSON(t *testing.T) {
	_ = Default
	_ = Fastest
}

func Test_parseVal(t *testing.T) {
	t.Run("Object", func(t *testing.T) {
		var v Value
		const input = `{"foo":{"bar":1,"baz":[1,2,3.14],"200":null}}`
		i := ParseString(Default, input)
		assert.NoError(t, parseVal(i, &v))
		assert.Equal(t, `{foo: {bar: 1, baz: [1, 2, f3.14], 200: null}}`, v.String())

		buf := new(bytes.Buffer)
		s := NewStream(Default, buf, 1024)
		v.Write(s)
		require.NoError(t, s.Flush())
		require.Equal(t, input, buf.String(), "encoded value should equal to input")
	})
	t.Run("Inputs", func(t *testing.T) {
		for _, tt := range []struct {
			Input string
		}{
			{Input: "1"},
			{Input: "0.0"},
		} {
			t.Run(tt.Input, func(t *testing.T) {
				var v Value
				input := []byte(tt.Input)
				i := ParseBytes(Default, input)
				require.NoError(t, parseVal(i, &v))

				buf := new(bytes.Buffer)
				s := NewStream(Default, buf, 1024)
				v.Write(s)
				require.NoError(t, s.Flush())
				require.Equal(t, tt.Input, buf.String(), "encoded value should equal to input")

				var otherValue Value
				i.ResetBytes(buf.Bytes())

				if err := parseVal(i, &otherValue); err != nil {
					t.Error(err)
					t.Log(hexEnc.Dump(input))
					t.Log(hexEnc.Dump(buf.Bytes()))
				}
			})
		}
	})
}

type ValType byte

const (
	ValInvalid ValType = iota
	ValStr
	ValInt
	ValFloat
	ValNull
	ValObj
	ValArr
	ValBool
)

// Value represents any json value as sum type.
type Value struct {
	Type   ValType
	Str    string
	Int    int64
	Float  float64
	Key    string
	KeySet bool // Key can be ""
	Bool   bool
	Child  []Value
}

// Write json representation of Value to Stream.
func (v Value) Write(s *Stream) {
	if v.KeySet {
		s.ObjField(v.Key)
	}
	switch v.Type {
	case ValStr:
		s.Str(v.Str)
	case ValFloat:
		s.WriteFloat64(v.Float)
	case ValInt:
		s.WriteInt64(v.Int)
	case ValBool:
		s.Bool(v.Bool)
	case ValNull:
		s.Null()
	case ValArr:
		s.ArrStart()
		for i, c := range v.Child {
			if i != 0 {
				s.More()
			}
			c.Write(s)
		}
		s.ArrEnd()
	case ValObj:
		s.ObjStart()
		for i, c := range v.Child {
			if i != 0 {
				s.More()
			}
			c.Write(s)
		}
		s.ObjEnd()
	default:
		panic(v.Type)
	}
}

func (v Value) String() string {
	var b strings.Builder
	if v.KeySet {
		if v.Key == "" {
			b.WriteString("<blank>")
		}
		b.WriteString(v.Key)
		b.WriteString(": ")
	}
	switch v.Type {
	case ValStr:
		b.WriteString(`"` + v.Str + `"'`)
	case ValFloat:
		b.WriteRune('f')
		b.WriteString(fmt.Sprintf("%v", v.Float))
	case ValInt:
		b.WriteString(strconv.FormatInt(v.Int, 10))
	case ValBool:
		b.WriteString(strconv.FormatBool(v.Bool))
	case ValNull:
		b.WriteString("null")
	case ValArr:
		b.WriteString("[")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("]")
	case ValObj:
		b.WriteString("{")
		for i, c := range v.Child {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.String())
		}
		b.WriteString("}")
	default:
		panic(v.Type)
	}
	return b.String()
}

func parseVal(i *Iter, v *Value) error {
	switch i.Next() {
	case Invalid:
		return xerrors.New("invalid")
	case Number:
		n, err := i.Number()
		if err != nil {
			return xerrors.Errorf("number: %w", err)
		}
		idx := strings.Index(n.String(), ".")
		if (idx > 0 && idx != len(n.String())-1) || strings.Contains(n.String(), "e") {
			f, err := n.Float64()
			if err != nil {
				return xerrors.Errorf("float: %w", err)
			}
			v.Float = f
			v.Type = ValFloat
		} else {
			f, err := n.Int64()
			if err != nil {
				return xerrors.Errorf("int: %w", err)
			}
			v.Int = f
			v.Type = ValInt
		}
	case String:
		s, err := i.Str()
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		v.Str = s
		v.Type = ValStr
	case Nil:
		if err := i.Null(); err != nil {
			return xerrors.Errorf("null: %w", err)
		}
		v.Type = ValNull
	case Bool:
		b, err := i.Bool()
		if err != nil {
			return xerrors.Errorf("bool: %w", err)
		}
		v.Bool = b
		v.Type = ValBool
	case Object:
		v.Type = ValObj
		if err := i.Object(func(i *Iter, s string) error {
			var elem Value
			if err := parseVal(i, &elem); err != nil {
				return xerrors.Errorf("elem: %w", err)
			}
			elem.Key = s
			elem.KeySet = true
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return xerrors.Errorf("obj: %w", err)
		}
		return nil
	case Array:
		v.Type = ValArr
		if err := i.Array(func(i *Iter) error {
			var elem Value
			if err := parseVal(i, &elem); err != nil {
				return xerrors.Errorf("elem: %w", err)
			}
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return xerrors.Errorf("array: %w", err)
		}
		return nil
	default:
		panic(i.Next())
	}
	return nil
}

// requireCompat fails if `encoding/json` will encode v differently than exp.
func requireCompat(t testing.TB, exp []byte, v interface{}) {
	t.Helper()
	buf, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, exp, buf)
}
