package bench

import (
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/simdjson-go"
)

// HelloWorld case.
//
// Example:
//	{"message": "Hello, world!"}

const (
	helloWorldField   = "message"
	helloWorldMessage = "Hello, world!"
	helloWorld        = `{"message": "Hello, world!"}`
)

//easyjson:json
type HelloWorld struct {
	Message string `json:"message"`
}

func (w HelloWorld) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart(helloWorldField)
	e.Str(w.Message)
	e.ObjEnd()
}

func (w HelloWorld) Write(wr *jx.Writer) {
	wr.ObjStart()
	wr.RawStr(`"message":`)
	wr.Str(w.Message)
	wr.ObjEnd()
}

func (w *HelloWorld) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case helloWorldField:
			v, err := d.Str()
			if err != nil {
				return err
			}
			w.Message = v
			return nil
		default:
			return d.Skip()
		}
	})
}

func (w *HelloWorld) DecodeSIMD(data []byte, reuse *simdjson.ParsedJson) (*simdjson.ParsedJson, error) {
	pj, err := simdjson.Parse(data, reuse, simdjson.WithCopyStrings(false))
	if err != nil {
		return nil, err
	}
	if err := pj.ForEach(func(i simdjson.Iter) error {
		typ := i.Advance()
		switch typ {
		case simdjson.TypeString:
			v, err := i.String()
			if err != nil {
				return err
			}
			w.Message = v
			return nil
		default:
			return errors.New("unexpected type")
		}
	}); err != nil {
		return nil, err
	}

	return pj, nil
}

func (w HelloWorld) EncodeIter(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField(helloWorldField)
	s.WriteString(w.Message)
	s.WriteObjectEnd()
}
