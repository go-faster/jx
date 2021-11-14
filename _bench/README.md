# Benchmarks

* jx
* jsoniter (should same as jx in most cases)
* encoding/json
* bytedance/sonic
* mailru/easyjson
* pquerna/ffjson

## Hello world

```json
{"message": "Hello, world!"}
```
```
goos: linux
goarch: amd64
pkg: bench
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkHelloWorld/Encode/std-32            11915048  94.82 ns/op  295.29 MB/s   0 B/op   0 allocs/op
BenchmarkHelloWorld/Encode/sonic-32          14639572  84.21 ns/op  332.51 MB/s   21 B/op  1 allocs/op
BenchmarkHelloWorld/Encode/ffjson-32         13557745  81.35 ns/op  344.17 MB/s   16 B/op  1 allocs/op
BenchmarkHelloWorld/Encode/jx-32             31989187  33.22 ns/op  842.80 MB/s   0 B/op   0 allocs/op
BenchmarkHelloWorld/Encode/json-iterator-32  32093022  32.73 ns/op  855.53 MB/s   0 B/op   0 allocs/op
BenchmarkHelloWorld/Encode/easyjson-32       47746503  21.69 ns/op  1290.65 MB/s  0 B/op   0 allocs/op
```
