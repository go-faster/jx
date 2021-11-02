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
	e.ObjStart()         // {
	e.ObjField("values") // "values":
	e.ArrStart()         // [
	for i, v := range []int{4, 8, 15, 16, 23, 42} {
		if i != 0 {
			e.More() // ,
		}
		e.Int(v)
	}
	e.ArrEnd() // ]
	e.ObjEnd() // }
	fmt.Println(e)
	fmt.Println("Buffer len:", len(e.Bytes()))
	// Output: {"values":[4,8,15,16,23,42]}
	// Buffer len: 28
}

func ExampleValid() {
	fmt.Println(jx.Valid([]byte(`{"field": "value"}`)))
	fmt.Println(jx.Valid([]byte(`"Hello, world!"`)))
	fmt.Println(jx.Valid([]byte(`["foo"}`)))
	// Output: true
	// true
	// false
}

func ExampleDecoder_Capture() {
	d := jx.DecodeStr(`["foo", "bar", "baz"]`)
	var elems int
	// NB: Currently Capture does not support io.Reader, only buffers.
	if err := d.Capture(func(d *jx.Decoder) error {
		// Everything decoded in this callback will be rolled back.
		return d.Arr(func(d *jx.Decoder) error {
			elems++
			return d.Skip()
		})
	}); err != nil {
		panic(err)
	}
	// Decoder is rolled back to state before "Capture" call.
	fmt.Println("Read", elems, "elements on first pass")
	fmt.Println("Next element is", d.Next(), "again")

	// Output:
	// Read 3 elements on first pass
	// Next element is array again
}
