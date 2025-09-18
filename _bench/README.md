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
$ GOEXPERIMENT=jsonv2 go test -bench .
goos: linux
goarch: amd64
pkg: bench
cpu: AMD Ryzen 9 5950X 16-Core Processor
```

## HelloWorld

### Encode

| Benchmark        | ns/op       | MB/s         | B/op     | allocs/op   |
|------------------|-------------|--------------|----------|-------------|
| jx/Encoder-32    | 38.51 ns/op | 727.07 MB/s  | 0 B/op   | 0 allocs/op |
| jx/Writer-32     | 19.38 ns/op | 1444.44 MB/s | 0 B/op   | 0 allocs/op |
| std-32           | 146.6 ns/op | 190.97 MB/s  | 0 B/op   | 0 allocs/op |
| sonnet-32        | 299.0 ns/op | 93.65 MB/s   | 288 B/op | 4 allocs/op |
| json-iterator-32 | 32.10 ns/op | 872.15 MB/s  | 0 B/op   | 0 allocs/op |
| easyjson-32      | 19.32 ns/op | 1449.63 MB/s | 0 B/op   | 0 allocs/op |
| ffjson-32        | 48.90 ns/op | 572.65 MB/s  | 16 B/op  | 1 allocs/op |
| stdv2-32         | 137.4 ns/op | 203.82 MB/s  | 0 B/op   | 0 allocs/op |
| Baseline-32      | 3.414 ns/op | 8200.33 MB/s | 0 B/op   | 0 allocs/op |

### Scan

| Benchmark   | ns/op       | MB/s         | B/op    | allocs/op   |
|-------------|-------------|--------------|---------|-------------|
| jx-32       | 26.39 ns/op | 1061.05 MB/s | 0 B/op  | 0 allocs/op |
| jscan-32    | 39.98 ns/op | 700.41 MB/s  | 0 B/op  | 0 allocs/op |
| simdjson-32 | 177.9 ns/op | 157.40 MB/s  | 16 B/op | 1 allocs/op |
| fastjson-32 | 29.50 ns/op | 949.10 MB/s  | 0 B/op  | 0 allocs/op |
### Decode


| Benchmark   | ns/op       | MB/s        | B/op    | allocs/op   |
|-------------|-------------|-------------|---------|-------------|
| jx-32       | 54.13 ns/op | 517.29 MB/s | 16 B/op | 1 allocs/op |
| simdjson-32 | 243.3 ns/op | 115.08 MB/s | 24 B/op | 2 allocs/op |
| std-32      | 235.2 ns/op | 119.06 MB/s | 0 B/op  | 0 allocs/op |
| fastjson-32 | 81.31 ns/op | 344.37 MB/s | 24 B/op | 2 allocs/op |
| stdv2-32    | 158.0 ns/op | 177.22 MB/s | 0 B/op  | 0 allocs/op |

## Small

### Encode
| Benchmark     | ns/op       | MB/s        | B/op     | allocs/op   |
|---------------|-------------|-------------|----------|-------------|
| jx/Encoder-32 | 900.1 ns/op | 379.95 MB/s | 0 B/op   | 0 allocs/op |
| jx/Writer-32  | 479.4 ns/op | 713.45 MB/s | 0 B/op   | 0 allocs/op |
| std-32        | 1344 ns/op  | 254.46 MB/s | 0 B/op   | 0 allocs/op |
| sonnet-32     | 940.1 ns/op | 363.80 MB/s | 288 B/op | 4 allocs/op |
| easyjson-32   | 431.6 ns/op | 792.36 MB/s | 0 B/op   | 0 allocs/op |

### Decode
| Benchmark   | ns/op      | MB/s        | B/op     | allocs/op    |
|-------------|------------|-------------|----------|--------------|
| EasyJet-32  | 1533 ns/op | 223.12 MB/s | 544 B/op | 14 allocs/op |
| std-32      | 3403 ns/op | 100.51 MB/s | 392 B/op | 15 allocs/op |
| sonnet-32   | 1710 ns/op | 199.95 MB/s | 476 B/op | 20 allocs/op |
| jx-32       | 1215 ns/op | 281.42 MB/s | 416 B/op | 20 allocs/op |
| fastjson-32 | 1378 ns/op | 248.20 MB/s | 416 B/op | 20 allocs/op |
| stdv2-32    | 2970 ns/op | 115.16 MB/s | 392 B/op | 15 allocs/op |

### Scan
| Benchmark   | ns/op       | MB/s        | B/op    | allocs/op   |
|-------------|-------------|-------------|---------|-------------|
| jx-32       | 549.5 ns/op | 622.35 MB/s | 0 B/op  | 0 allocs/op |
| jscan-32    | 578.4 ns/op | 591.24 MB/s | 0 B/op  | 0 allocs/op |
| simdjson-32 | 1041 ns/op  | 328.42 MB/s | 16 B/op | 1 allocs/op |
| fastjson-32 | 540.0 ns/op | 633.35 MB/s | 0 B/op  | 0 allocs/op |
