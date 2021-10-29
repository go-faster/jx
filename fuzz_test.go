//go:build go1.18

package jx

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzValid(f *testing.F) {
	f.Add("{}")
	f.Add(`{"foo": "bar"}`)
	f.Add(``)
	f.Add(`"foo"`)
	f.Add(`"{"`)
	f.Add(`"{}"`)
	f.Fuzz(func(t *testing.T, queryStr string) {
		Valid([]byte(queryStr))
	})
}

func FuzzDecEnc(f *testing.F) {
	f.Add([]byte("{}"))
	f.Add([]byte(`"foo"`))
	f.Add([]byte(`123"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"foo": {"bar": 1, "baz": [1, 2, 3]}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		i := Default.GetIter(nil)
		i.ResetBytes(data)
		defer Default.PutIter(i)

		// Parsing to v.
		var v Value
		if parseVal(i, &v) != nil {
			t.Skip()
		}
		if v.Type == ValInvalid {
			t.Skip()
		}
		// Writing v to buf.
		var buf bytes.Buffer
		s := Default.GetStream(&buf)
		v.Write(s)
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}

		// Parsing from buf to new value.
		i.ResetBytes(buf.Bytes())
		var parsed Value
		if err := parseVal(i, &parsed); err != nil {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s\nData:  %s",
				err, buf.Bytes(), v, data)
		}
		if !reflect.DeepEqual(parsed, v) {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s != %s \nData:  %s",
				nil, buf.Bytes(), parsed, v, data)
		}
		// Writing parsed value to newBuf.
		var newBuf bytes.Buffer
		s.Reset(&newBuf)
		parsed.Write(s)
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(newBuf.Bytes(), buf.Bytes()) {
			t.Fatalf("%s != %s", &newBuf, &buf)
		}
	})
}
