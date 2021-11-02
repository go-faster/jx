# jx [![Go Reference](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/ogen-go/jx#section-documentation) [![codecov](https://img.shields.io/codecov/c/github/ogen-go/jx?label=cover)](https://codecov.io/gh/ogen-go/jx)

Package jx implements encoding and decoding of json.
Lightweight fork of [jsoniter](https://github.com/json-iterator/go).

```console
go get github.com/ogen-go/jx
```

## Features
* Directly encode and decode json values
* No reflect or `interface{}`
* Pools and direct buffer access for less (or none) allocations
* Multi-pass decoding
* Validation

See [usage](#Usage) for examples. Mostly suitable for low-level json manipulation
when high performance and control are needed. Used in [ogen](https://github.com/ogen-go/ogen) for
json (un)marshaling code generation.

## Why

Most of [jsoniter](https://github.com/json-iterator/go) issues are caused by necessity
to be drop-in replacement for standard `encoding/json`. Removing such constrains greatly
simplified implementation and reduced scope, allowing to focus on json stream processing.

* Reduced scope
  * No reflection
  * No `encoding/json` adapter
  * 2.2K SLOC vs 8.5K in `jsoniter`
* Fuzzing, improved test coverage
* Drastically refactored and simplified
  * Explicit error returns
  * No `Config` or `API`


## Usage

### Decode

Use [jx.Decoder](https://pkg.go.dev/github.com/ogen-go/jx#Decoder). Zero value is valid,
but constructors are available for convenience:
  * [jx.Decode(reader io.Reader, bufSize int)](https://pkg.go.dev/github.com/ogen-go/jx#Decode) for `io.Reader`
  * [jx.DecodeBytes([]byte)](https://pkg.go.dev/github.com/ogen-go/jx#Decode)  for byte slices
  * [jx.DecodeStr(string)](https://pkg.go.dev/github.com/ogen-go/jx#Decode) for strings

To reuse decoders and their buffers, use [jx.GetDecoder](https://pkg.go.dev/github.com/ogen-go/jx#GetDecoder)
and [jx.PutDecoder](https://pkg.go.dev/github.com/ogen-go/jx#PutDecoder) alongside with reset functions:
* [jx.Decoder.Reset(io.Reader)](https://pkg.go.dev/github.com/ogen-go/jx#Decoder.Reset) to reset to new `io.Reader`
* [jx.Decoder.ResetBytes([]byte)](https://pkg.go.dev/github.com/ogen-go/jx#Decoder.ResetBytes) to decode another byte slice

Decoder is reset on `PutDecoder`.

```go
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
```

### Encode
Use [jx.Encoder](https://pkg.go.dev/github.com/ogen-go/jx#Encoder). Zero value is valid, reuse with
[jx.GetEncoder](https://pkg.go.dev/github.com/ogen-go/jx#GetEncoder),
[jx.PutEncoder](https://pkg.go.dev/github.com/ogen-go/jx#PutEncoder) and
[jx.Encoder.Reset()](https://pkg.go.dev/github.com/ogen-go/jx#Encoder.Reset). Encoder is reset on `PutEncoder`.
```go
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
```

## Validate

Check that byte slice is valid json with [jx.Valid](https://pkg.go.dev/github.com/ogen-go/jx#Valid):

```go
fmt.Println(jx.Valid([]byte(`{"field": "value"}`))) // true
fmt.Println(jx.Valid([]byte(`"Hello, world!"`)))    // true
fmt.Println(jx.Valid([]byte(`["foo"}`)))            // false
```

## Capture
The [jx.Decoder.Capture](https://pkg.go.dev/github.com/ogen-go/jx#Decoder.Capture) method allows to unread everything is read in callback.
Useful for multi-pass parsing:
```go
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
```

## ObjBytes

The `Decoder.ObjBytes` method tries not to allocate memory for keys, reusing existing buffer.
```go
d := DecodeStr(`{"id":1,"randomNumber":10}`)
d.ObjBytes(func(d *Decoder, key []byte) error {
    switch string(key) {
    case "id":
    case "randomNumber":
    }
    return d.Skip()
})
```

# Roadmap
- [ ] Rework `json.Number`
- [ ] Rework `Any`
- [ ] Support `Raw` decoding

# License
MIT, same as jsoniter
