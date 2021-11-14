# Benchmarks

* jx
* jsoniter (should same as jx in most cases)
* encoding/json
* bytedance/sonic

## Hello world

```json
{"message": "Hello, world!"}
```
```
goos: linux
goarch: amd64
pkg: bench
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkHelloWorld
BenchmarkHelloWorld/Encode
BenchmarkHelloWorld/Encode/jx
BenchmarkHelloWorld/Encode/jx-32      	31651273	        32.37 ns/op	 864.90 MB/s	       0 B/op	       0 allocs/op
BenchmarkHelloWorld/Encode/std
BenchmarkHelloWorld/Encode/std-32     	11931632	        98.73 ns/op	 283.59 MB/s	       0 B/op	       0 allocs/op
BenchmarkHelloWorld/Encode/sonic
BenchmarkHelloWorld/Encode/sonic-32   	15214280	        81.68 ns/op	 342.79 MB/s	      22 B/op	       1 allocs/op
BenchmarkHelloWorld/Encode/json-iterator
BenchmarkHelloWorld/Encode/json-iterator-32         	37347980	        32.62 ns/op	 858.48 MB/s	       0 B/op	       0 allocs/op
PASS
```
