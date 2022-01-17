//go:build !gofuzz && go1.17
// +build !gofuzz,go1.17

package jx

import (
	"embed"
	"path"
	"testing"

	"github.com/go-faster/errors"
)

var (
	//go:embed testdata/medium.json
	benchData []byte
	//go:embed testdata
	testdata embed.FS
)

func runTestdata(fatal func(...interface{}), cb func(name string, data []byte)) {
	dir, err := testdata.ReadDir("testdata")
	if err != nil {
		fatal(err)
	}
	for _, e := range dir {
		if e.IsDir() {
			continue
		}
		runTestdataFile(e.Name(), fatal, cb)
	}
}

func runTestdataFile(file string, fatal func(...interface{}), cb func(name string, data []byte)) {
	name := path.Join("testdata", file)
	data, err := testdata.ReadFile(name)
	if err != nil {
		fatal(err)
	}
	cb(file, data)
}

func BenchmarkFile_Decode(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchData)))
	d := Decode(nil, 4096)

	for n := 0; n < b.N; n++ {
		d.ResetBytes(benchData)
		if err := d.Arr(func(d *Decoder) error {
			return d.ObjBytes(func(d *Decoder, key []byte) error {
				switch string(key) {
				case "person", "company": // ok
				default:
					return errors.New("unexpected key")
				}
				switch d.Next() {
				case Object:
					return d.ObjBytes(func(d *Decoder, key []byte) error {
						switch d.Next() {
						case String:
							_, err := d.StrBytes()
							return err
						case Number:
							_, err := d.Num()
							return err
						case Null:
							return d.Null()
						default:
							return d.Skip()
						}
					})
				default:
					return d.Skip()
				}
			})
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValid(b *testing.B) {
	d := GetDecoder()
	runTestdata(b.Fatal, func(name string, data []byte) {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.ResetBytes(data)
				if err := d.Validate(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

func encodeSmallObject(e *Encoder) {
	e.ObjStart()
	e.FieldStart("data_array")
	e.ArrStart()
	e.Int(5467889)
	e.Int(456717)
	e.Int(5789935)
	e.ArrEnd()
	e.ObjEnd()
}

func BenchmarkEncoder_ObjStart(b *testing.B) {
	e := GetEncoder()
	encodeSmallObject(e)
	setBytes(b, e)
	if e.String() != `{"data_array":[5467889,456717,5789935]}` {
		b.Fatal(e)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e.Reset()
		encodeSmallObject(e)
	}
}

func encodeSmallCallback(e *Encoder) {
	e.Obj(func(e *Encoder) {
		e.Field("foo", func(e *Encoder) {
			e.Arr(func(e *Encoder) {
				e.Int(100)
				e.Int(200)
				e.Int(300)
			})
		})
	})
}

func setBytes(b *testing.B, e *Encoder) {
	b.Helper()
	b.SetBytes(int64(len(e.Bytes())))
}

func BenchmarkEncoder_Obj(b *testing.B) {
	e := GetEncoder()
	b.ReportAllocs()

	encodeSmallCallback(e)
	setBytes(b, e)
	if string(e.Bytes()) != `{"foo":[100,200,300]}` {
		b.Fatal("mismatch")
	}

	for i := 0; i < b.N; i++ {
		e.Reset()
		encodeSmallCallback(e)
	}
}
