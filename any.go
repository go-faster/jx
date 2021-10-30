package jx

import (
	"fmt"
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

// Any reads Any value from r.
func (r *Reader) Any() (Any, error) {
	var v Any
	if err := v.Read(r); err != nil {
		return Any{}, err
	}
	return v, nil
}

// Any writes Any value to w.
func (w *Writer) Any(a Any) error {
	return a.Write(w)
}

func (v *Any) Read(r *Reader) error {
	switch r.Next() {
	case Invalid:
		return xerrors.New("invalid")
	case Number:
		n, err := r.Number()
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
		s, err := r.Str()
		if err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		v.Str = s
		v.Type = AnyStr
	case Nil:
		if err := r.Null(); err != nil {
			return xerrors.Errorf("null: %w", err)
		}
		v.Type = AnyNull
	case Bool:
		b, err := r.Bool()
		if err != nil {
			return xerrors.Errorf("bool: %w", err)
		}
		v.Bool = b
		v.Type = AnyBool
	case Object:
		v.Type = AnyObj
		if err := r.Obj(func(r *Reader, s string) error {
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
		if err := r.Array(func(r *Reader) error {
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
		return xerrors.Errorf("unexpected type %s", r.Next())
	}
	return nil
}

// Write json representation of Any to Writer.
func (v Any) Write(w *Writer) error {
	if v.KeyValid {
		w.ObjField(v.Key)
	}
	switch v.Type {
	case AnyStr:
		w.Str(v.Str)
	case AnyFloat:
		if err := w.Float64(v.Float); err != nil {
			return err
		}
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
			if err := c.Write(w); err != nil {
				return err
			}
		}
		w.ArrEnd()
	case AnyObj:
		w.ObjStart()
		for i, c := range v.Child {
			if i != 0 {
				w.More()
			}
			if err := c.Write(w); err != nil {
				return err
			}
		}
		w.ObjEnd()
	default:
		return xerrors.Errorf("unexpected type %d", v.Type)
	}
	return nil
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
		b.WriteString(`"` + v.Str + `"'`)
	case AnyFloat:
		b.WriteRune('f')
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
