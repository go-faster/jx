# jx [![Go Reference](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/ogen-go/jx#section-documentation) [![codecov](https://img.shields.io/codecov/c/github/ogen-go/jx?label=cover)](https://codecov.io/gh/ogen-go/jx)

Fast json streaming in go. Buffered encoding and decoding of json values.
Lightweight fork of [jsoniter](https://github.com/json-iterator/go).

```console
go get github.com/ogen-go/jx
```

## Features
* Reduced scope
  * No reflection
  * No `encoding/json` adapter
  * 2.2K SLOC vs 8.5K in `jsoniter`
* Fuzzing, improved test coverage
* Drastically refactored and simplified
  * Explicit error returns
  * No `Config` or `API`

## Why

Most of [jsoniter](https://github.com/json-iterator/go) issues are caused by necessity
to be drop-in replacement for standard `encoding/json`. Removing such constrains greatly
simplified implementation and reduced scope, allowing to focus on json stream processing.

## Capture

The `Decoder.Capture` method allows to unread everything is read in callback.
This is useful for multi-pass parsing:
```go
func TestDecoder_Capture(t *testing.T) {
	d := DecodeString(`["foo", "bar", "baz"]`)
	var elems int
	if err := d.Capture(func(d *Decoder) error {
		return d.Arr(func(d *Decoder) error {
			elems++
			return d.Skip()
		})
	}); err != nil {
		t.Fatal(err)
	}
	// Buffer is rolled back to state before "Capture" call:
	require.Equal(t, Array, d.Next())
	require.Equal(t, 3, elems)
}
```

## ObjectBytes

The `Decoder.ObjectBytes` method tries not to allocate memory for keys, reusing existing buffer:
```go
d := DecodeString(`{"id":1,"randomNumber":10}`)
d.ObjectBytes(func(d *Decoder, key []byte) error {
    switch string(key) {
    case "id":
    case "randomNumber":
    }
    return d.Skip()
})
```
