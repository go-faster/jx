# Benchmarks

* jx
* jsoniter (should same as jx in most cases)
* encoding/json
* bytedance/sonic
* mailru/easyjson
* pquerna/ffjson

| name                            | ns/op | MB/s    | B/op | Allocs |
|---------------------------------|-------|---------|------|--------|
| HelloWorld/Encode/Baseline      | 5.238 | 5345.77 | 0    | 0      |
| HelloWorld/Encode/easyjson      | 28.99 | 965.97  | 0    | 0      |
| HelloWorld/Encode/ffjson        | 93.15 | 300.59  | 16   | 1      |
| HelloWorld/Encode/json-iterator | 44.71 | 626.24  | 0    | 0      |
| HelloWorld/Encode/jx/Encoder    | 34.89 | 802.51  | 0    | 0      |
| HelloWorld/Encode/jx/Writer     | 18.38 | 1523.38 | 0    | 0      |
| HelloWorld/Encode/sonic         | 104.4 | 268.32  | 21   | 1      |
| HelloWorld/Encode/std           | 101   | 277.3   | 0    | 0      |
| HelloWorld/Scan/jscan           | 57.41 | 487.75  | 0    | 0      |
| HelloWorld/Scan/jx              | 55.86 | 501.29  | 0    | 0      |
| Small/Decode/easyjson           | 1629  | 208.1   | 32   | 7      |
| Small/Decode/sonic              | 982.5 | 345.06  | 0    | 0      |
| Small/Decode/std                | 8789  | 38.57   | 400  | 24     |
| Small/Encode/easyjson           | 647.2 | 523.76  | 0    | 0      |
| Small/Encode/jx/Encoder         | 769.1 | 440.76  | 0    | 0      |
| Small/Encode/jx/Writer          | 408.4 | 829.98  | 0    | 0      |
| Small/Encode/sonic              | 532.6 | 636.54  | 32   | 1      |
| Small/Encode/std                | 1144  | 296.26  | 0    | 0      |
| Small/Scan/jscan                | 802   | 422.69  | 0    | 0      |
| Small/Scan/jx                   | 1118  | 303.2   | 0    | 0      |
