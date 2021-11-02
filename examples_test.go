package jx_test

import (
	"fmt"

	"github.com/ogen-go/jx"
)

func ExampleDecodeStr() {
	d := jx.DecodeStr(`{"values":[4,8,15,16,23,42]}`)

	// Save all integers from "values" array to slice.
	var values []int

	// Iterate over each object field.
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		switch key {
		case "values":
			// Iterate over each array element.
			return d.Arr(func(d *jx.Decoder) error {
				v, err := d.Int()
				if err != nil {
					return err
				}
				values = append(values, v)
				return nil
			})
		default:
			// Skip unknown fields if any.
			return d.Skip()
		}
	}); err != nil {
		panic(err)
	}

	fmt.Println(values)
	// Output: [4 8 15 16 23 42]
}

func ExampleEncoder_String() {
	var e jx.Encoder
	e.ObjStart()
	e.ObjField("values")
	e.ArrStart()
	for i, v := range []int{4, 8, 15, 16, 23, 42} {
		if i != 0 {
			e.More()
		}
		e.Int(v)
	}
	e.ArrEnd()
	e.ObjEnd()
	fmt.Println(e)
	// Output: {"values":[4,8,15,16,23,42]}
}
