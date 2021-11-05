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

func ExampleDecoder_Raw() {
	d := jx.DecodeStr(`{"foo": [1, 2, 3]}`)

	var raw jx.Raw
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		v, err := d.Raw()
		if err != nil {
			return err
		}
		raw = v
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println(raw.Type(), raw)
	// Output:
	// array [1, 2, 3]
}

func ExampleDecoder_Num() {
	// Can decode numbers and number strings.
	d := jx.DecodeStr(`{"foo": "10531.0"}`)

	var n jx.Num
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		v, err := d.Num()
		if err != nil {
			return err
		}
		n = v
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println(n)
	fmt.Println("positive:", n.Positive())

	// Can decode floats with zero fractional part as integers:
	v, err := n.Int64()
	if err != nil {
		panic(err)
	}
	fmt.Println("int64:", v)
	// Output:
	// "10531.0"
	// positive: true
	// int64: 10531
}
