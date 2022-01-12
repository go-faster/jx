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
BenchmarkSmall/Encode/jx/Encoder-32  1977728  617.9 ns/op   548.60 MB/s  0 B/op   0 allocs/op
BenchmarkSmall/Encode/jx/Writer-32   4067817  294.9 ns/op  1149.71 MB/s  0 B/op   0 allocs/op
BenchmarkSmall/Encode/std-32         1254765  931.3 ns/op   363.99 MB/s  0 B/op   0 allocs/op
BenchmarkSmall/Encode/sonic-32       2534568  458.9 ns/op   738.80 MB/s  16 B/op  1 allocs/op
BenchmarkSmall/Encode/easyjson-32    2493712  469.3 ns/op   722.29 MB/s  0 B/op   0 allocs/op
PASS
ok      bench   8.411s
```
