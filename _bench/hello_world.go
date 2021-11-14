package bench

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/go-faster/jx"
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

func (w HelloWorld) EncodeIter(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField(helloWorldField)
	s.WriteString(w.Message)
	s.WriteObjectEnd()
}
