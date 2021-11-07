//go:build !gofuzz && go1.17
// +build !gofuzz,go1.17

package jx

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/go-faster/errors"
)

//go:embed testdata/file.json
var benchData []byte

func Benchmark_large_file(b *testing.B) {
	b.Run("JX", func(b *testing.B) {
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
	})
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchData)))

		type T struct {
			Person struct {
				ID   string `json:"id"`
				Name struct {
					FullName   string `json:"fullName"`
					GivenName  string `json:"givenName"`
					FamilyName string `json:"familyName"`
				} `json:"name"`
				Email    string `json:"email"`
				Gender   string `json:"gender"`
				Location string `json:"location"`
				Geo      struct {
					City    string  `json:"city"`
					State   string  `json:"state"`
					Country string  `json:"country"`
					Lat     float64 `json:"lat"`
					Lng     float64 `json:"lng"`
				} `json:"geo"`
				Bio        string `json:"bio"`
				Site       string `json:"site"`
				Avatar     string `json:"avatar"`
				Employment struct {
					Name   string `json:"name"`
					Title  string `json:"title"`
					Domain string `json:"domain"`
				} `json:"employment"`
				Facebook struct {
					Handle string `json:"handle"`
				} `json:"facebook"`
				Github struct {
					Handle    string `json:"handle"`
					ID        int    `json:"id"`
					Avatar    string `json:"avatar"`
					Company   string `json:"company"`
					Blog      string `json:"blog"`
					Followers int    `json:"followers"`
					Following int    `json:"following"`
				} `json:"github"`
				Twitter struct {
					Handle    string          `json:"handle"`
					ID        int             `json:"id"`
					Bio       json.RawMessage `json:"bio"`
					Followers int             `json:"followers"`
					Following int             `json:"following"`
					Statuses  int             `json:"statuses"`
					Favorites int             `json:"favorites"`
					Location  string          `json:"location"`
					Site      string          `json:"site"`
					Avatar    json.RawMessage `json:"avatar"`
				} `json:"twitter"`
				Linkedin struct {
					Handle string `json:"handle"`
				} `json:"linkedin"`
				Googleplus struct {
					Handle json.RawMessage `json:"handle"`
				} `json:"googleplus"`
				Angellist struct {
					Handle    string `json:"handle"`
					ID        int    `json:"id"`
					Bio       string `json:"bio"`
					Blog      string `json:"blog"`
					Site      string `json:"site"`
					Followers int    `json:"followers"`
					Avatar    string `json:"avatar"`
				} `json:"angellist"`
				Klout struct {
					Handle json.RawMessage `json:"handle"`
					Score  json.RawMessage `json:"score"`
				} `json:"klout"`
				Foursquare struct {
					Handle json.RawMessage `json:"handle"`
				} `json:"foursquare"`
				Aboutme struct {
					Handle string          `json:"handle"`
					Bio    json.RawMessage `json:"bio"`
					Avatar json.RawMessage `json:"avatar"`
				} `json:"aboutme"`
				Gravatar struct {
					Handle  string          `json:"handle"`
					Urls    json.RawMessage `json:"urls"`
					Avatar  string          `json:"avatar"`
					Avatars []struct {
						URL  string `json:"url"`
						Type string `json:"type"`
					} `json:"avatars"`
				} `json:"gravatar"`
				Fuzzy bool `json:"fuzzy"`
			} `json:"person"`
			Company string `json:"company"`
		}

		buf := new(bytes.Reader)
		d := json.NewDecoder(buf)
		var target []T
		for n := 0; n < b.N; n++ {
			buf.Reset(benchData)
			if err := d.Decode(&target); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkValid(b *testing.B) {
	b.Run("JX", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchData)))
		var d Decoder
		for n := 0; n < b.N; n++ {
			d.ResetBytes(benchData)
			if err := d.Validate(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchData)))

		for n := 0; n < b.N; n++ {
			if !json.Valid(benchData) {
				b.Fatal("invalid")
			}
		}
	})
}

func Benchmark_std_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var result []struct{}
		err := json.Unmarshal(benchData, &result)
		if err != nil {
			b.Error(err)
		}
	}
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
