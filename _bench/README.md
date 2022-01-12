# Benchmarks

* jx
* jsoniter (should same as jx in most cases)
* encoding/json
* bytedance/sonic
* mailru/easyjson
* pquerna/ffjson

| name                            | ns/op | MB/s    | B/op | allocs/op |
|---------------------------------|-------|---------|------|-----------|
| HelloWorld/Encode/Baseline      | 3.75  | 7466.65 | 0    | 0         |
| HelloWorld/Encode/easyjson      | 21.46 | 1304.66 | 0    | 0         |
| HelloWorld/Encode/ffjson        | 105.8 | 264.53  | 16   | 1         |
| HelloWorld/Encode/json-iterator | 35.11 | 797.52  | 0    | 0         |
| HelloWorld/Encode/jx/Encoder    | 30.49 | 918.33  | 0    | 0         |
| HelloWorld/Encode/jx/Writer     | 14.27 | 1962.1  | 0    | 0         |
| HelloWorld/Encode/sonic         | 105.5 | 265.44  | 21   | 1         |
| HelloWorld/Encode/std           | 85.99 | 325.63  | 0    | 0         |
| HelloWorld/Scan/jscan           | 43.47 | 644.09  | 0    | 0         |
| HelloWorld/Scan/jx              | 45.93 | 609.63  | 0    | 0         |
| Small/Decode/easyjson           | 1381  | 245.46  | 32   | 7         |
| Small/Decode/sonic              | 743   | 456.27  | 1    | 0         |
| Small/Decode/std                | 7836  | 43.26   | 400  | 24        |
| Small/Encode/easyjson           | 455.2 | 744.76  | 0    | 0         |
| Small/Encode/jx/Encoder         | 626.9 | 540.75  | 0    | 0         |
| Small/Encode/jx/Writer          | 310.4 | 1092.08 | 0    | 0         |
| Small/Encode/sonic              | 454.6 | 745.74  | 18   | 1         |
| Small/Encode/std                | 871.4 | 389.03  | 0    | 0         |
| Small/Scan/jscan                | 631.2 | 537.04  | 0    | 0         |
| Small/Scan/jx                   | 849   | 399.29  | 0    | 0         |
