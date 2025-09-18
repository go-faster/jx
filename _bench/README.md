# Benchmarks

* jx
* jsoniter (should same as jx in most cases)
* fastjson
* encoding/json
* bytedance/sonic
* mailru/easyjson
* pquerna/ffjson
* simdjson-go

```
go test -bench .
goos: linux
goarch: amd64
pkg: bench
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkHelloWorld/Encode/jx/Encoder-32                29534337                38.32 ns/op      730.73 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Encode/jx/Writer-32                 57847146                19.71 ns/op     1420.67 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Encode/std-32                       14543918                77.42 ns/op      361.64 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Encode/sonnet-32                     3656274               321.0 ns/op        87.22 MB/s         288 B/op          4 allocs/op
BenchmarkHelloWorld/Encode/json-iterator-32             36615604                31.75 ns/op      881.92 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Encode/easyjson-32                  62660282                19.26 ns/op     1454.01 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Encode/ffjson-32                    22818333                51.58 ns/op      542.86 MB/s          16 B/op          1 allocs/op
BenchmarkHelloWorld/Encode/Baseline-32                  326291570                3.484 ns/op    8036.64 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Scan/jx-32                          43161399                26.95 ns/op     1038.82 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Scan/jscan-32                       30240622                40.31 ns/op      694.70 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Scan/simdjson-32                     6840558               176.4 ns/op       158.70 MB/s          16 B/op          1 allocs/op
BenchmarkHelloWorld/Scan/fastjson-32                    40265343                29.52 ns/op      948.42 MB/s           0 B/op          0 allocs/op
BenchmarkHelloWorld/Decode/jx-32                        24218804                46.91 ns/op      596.82 MB/s          16 B/op          1 allocs/op
BenchmarkHelloWorld/Decode/simdjson-32                   5015164               238.4 ns/op       117.44 MB/s          24 B/op          2 allocs/op
BenchmarkHelloWorld/Decode/std-32                        2747335               443.0 ns/op        63.20 MB/s         232 B/op          5 allocs/op
BenchmarkHelloWorld/Decode/fastjson-32                  13771676                86.02 ns/op      325.52 MB/s          24 B/op          2 allocs/op
BenchmarkSmall/Encode/jx/Encoder-32                      1314556               898.5 ns/op       380.64 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Encode/jx/Writer-32                       2332848               510.3 ns/op       670.22 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Encode/std-32                             1428535               835.0 ns/op       409.57 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Encode/sonnet-32                          1239702              1002 ns/op         341.32 MB/s         288 B/op          4 allocs/op
BenchmarkSmall/Encode/easyjson-32                        2621916               454.7 ns/op       752.08 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Decode/easyjson-32                         718976              1596 ns/op         214.27 MB/s         544 B/op         14 allocs/op
BenchmarkSmall/Decode/std-32                              206193              5187 ns/op          65.94 MB/s         712 B/op         27 allocs/op
BenchmarkSmall/Decode/sonnet-32                           590230              1822 ns/op         187.67 MB/s         476 B/op         20 allocs/op
BenchmarkSmall/Decode/jx-32                               915477              1236 ns/op         276.81 MB/s         416 B/op         20 allocs/op
BenchmarkSmall/Decode/fastjson-32                         885486              1139 ns/op         300.18 MB/s         216 B/op         13 allocs/op
BenchmarkSmall/Scan/jx-32                                2142826               555.2 ns/op       615.96 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Scan/jscan-32                             2065302               582.9 ns/op       586.71 MB/s           0 B/op          0 allocs/op
BenchmarkSmall/Scan/simdjson-32                          1000000              1040 ns/op         328.71 MB/s          16 B/op          1 allocs/op
BenchmarkSmall/Scan/fastjson-32                          2203672               541.7 ns/op       631.40 MB/s           0 B/op          0 allocs/op
PASS
ok      bench   43.505s
```
