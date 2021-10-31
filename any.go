package jx

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// AnyType is type of Any value.
type AnyType byte

// Possible types for Any.
const (
	AnyInvalid AnyType = iota
	AnyStr
	AnyInt
	AnyFloat
	AnyNull
	AnyObj
	AnyArr
	AnyBool
)

// Any represents any json value as sum type.
type Any struct {
	Type AnyType // zero value if AnyInvalid, can be AnyNull

	Str   string  // AnyStr
	Int   int64   // AnyInt
	Float float64 // AnyFloat
	Bool  bool    // AnyBool

	// Key in object. Valid only if KeyValid.
	Key string
	// KeyValid denotes whether Any is element of object.
	// Needed for representing Key that is blank.
	//
	// Can be true only for Child of AnyObj.
	KeyValid bool

	Child []Any // AnyArr or AnyObj
}

func floatInt(f float64, n int64) bool {
	if n == 0 && f == 0 {
		return true
	}
	i, frac := math.Modf(f)
	if frac != 0 {
		return false
	}
	return int64(i) == n
}

func (v Any) equalNumber(b Any) bool {
	if b.Type != AnyFloat && b.Type != AnyInt {
		return false
	}
	eq := v.Type == b.Type
	if eq && v.Type == AnyInt {
		return v.Int == b.Int
	}
	if eq && v.Type == AnyFloat {
		return v.Float == b.Float
	}
	if v.Type == AnyFloat {
		return floatInt(v.Float, b.Int)
	}
	return floatInt(b.Float, v.Int)
}

func (v Any) Equal(b Any) bool {
	if v.KeyValid && v.Key != b.Key {
		return false
	}
	if v.Type == AnyFloat || v.Type == AnyInt {
		return v.equalNumber(b)
	}
	if v.Type != b.Type {
		return false
	}
	switch v.Type {
	case AnyNull:
		return true
	case AnyBool:
		return v.Bool == b.Bool
	case AnyInvalid:
		return false
	case AnyStr:
		return v.Str == b.Str
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
		return xerrors.New("invalid")
	case Number:
		n, err := d.Number()
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
			v.Type = AnyFloat
		} else {
			f, err := n.Int64()
			if err != nil {
				return xerrors.Errorf("int: %w", err)
			}
			v.Int = f
			v.Type = AnyInt
		}
	case String:
		s, err := d.Str()
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		v.Str = s
		v.Type = AnyStr
	case Nil:
		if err := d.Null(); err != nil {
			return xerrors.Errorf("null: %w", err)
		}
		v.Type = AnyNull
	case Bool:
		b, err := d.Bool()
		if err != nil {
			return xerrors.Errorf("bool: %w", err)
		}
		v.Bool = b
		v.Type = AnyBool
	case Object:
		v.Type = AnyObj
		if err := d.Obj(func(r *Decoder, s string) error {
			var elem Any
			if err := elem.Read(r); err != nil {
				return xerrors.Errorf("elem: %w", err)
			}
			elem.Key = s
			elem.KeyValid = true
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return xerrors.Errorf("obj: %w", err)
		}
		return nil
	case Array:
		v.Type = AnyArr
		if err := d.Arr(func(r *Decoder) error {
			var elem Any
			if err := elem.Read(r); err != nil {
				return xerrors.Errorf("elem: %w", err)
			}
			v.Child = append(v.Child, elem)
			return nil
		}); err != nil {
			return xerrors.Errorf("array: %w", err)
		}
		return nil
	default:
		return xerrors.Errorf("unexpected type %s", d.Next())
	}
	return nil
}

// Write json representation of Any to Encoder.
func (v Any) Write(w *Encoder) {
	if v.KeyValid {
		w.ObjField(v.Key)
	}
	switch v.Type {
	case AnyStr:
		w.Str(v.Str)
	case AnyFloat:
		w.Float64(v.Float)
	case AnyInt:
		w.Int64(v.Int)
	case AnyBool:
		w.Bool(v.Bool)
	case AnyNull:
		w.Null()
	case AnyArr:
		w.ArrStart()
		for i, c := range v.Child {
			if i != 0 {
				w.More()
			}
			c.Write(w)
		}
		w.ArrEnd()
	case AnyObj:
		w.ObjStart()
		for i, c := range v.Child {
			if i != 0 {
				w.More()
			}
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
	case AnyFloat:
		b.WriteString(fmt.Sprintf("%v", v.Float))
	case AnyInt:
		b.WriteString(strconv.FormatInt(v.Int, 10))
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
		b.WriteString("<unknown>")
	}
	return b.String()
}
