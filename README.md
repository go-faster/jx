# jx [![Go Reference](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/ogen-go/jx#section-documentation) [![codecov](https://img.shields.io/codecov/c/github/ogen-go/jx?label=cover)](https://codecov.io/gh/ogen-go/jx)

Fast json for go. Lightweight fork of [jsoniter](https://github.com/json-iterator/go).

## Features
* Reduced scope (no reflection or `encoding/json` adapter)
* Fuzzing, improved test coverage
* Drastically refactored and simplified
  * Explicit error returns
  * No `Config` or `API`

## Capture

The `Iter.Capture` method allows to unread everything is read in callback.
This is useful for multi-pass parsing:
```go
func TestIter_Capture(t *testing.T) {
	i := ParseString(`["foo", "bar", "baz"]`)
	var elems int
	if err := i.Capture(func(i *Iter) error {
		return i.Array(func(i *Iter) error {
			elems++
			return i.Skip()
		})
	}); err != nil {
		t.Fatal(err)
	}
	// Buffer is rolled back to state before "Capture" call:
	require.Equal(t, Array, i.Next())
	require.Equal(t, 3, elems)
}
```

## ObjBytes

The `Iter.ObjBytes` method tries not to allocate memory for keys, reusing existing buffer:
```go
i := ParseString(`{"id":1,"randomNumber":10}`)
i.ObjBytes(func(i *Iter, key []byte) error {
    switch string(key) {
    case "id":
    case "randomNumber":
    }
    return i.Skip()
})
```
