//go:build go1.18

package jx

import (
	"bytes"
	"errors"
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

func FuzzValues(f *testing.F) {
	f.Add(int64(1), "hello")
	f.Add(int64(1534564316421), " привет ")
	f.Fuzz(func(t *testing.T, n int64, str string) {
		buf := new(bytes.Buffer)
		s := Default.GetStream(buf)
		defer Default.PutStream(s)

		s.ArrStart()
		s.WriteInt64(n)
		s.More()
		s.Str(str)
		s.ArrEnd()

		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}

		i := Default.GetIter(buf.Bytes())
		var (
			nGot int64
			sGot string
		)
		if err := i.Array(func(i *Iter) error {
			var err error
			switch i.Next() {
			case Number:
				nGot, err = i.Int64()
			case String:
				sGot, err = i.Str()
			default:
				err = errors.New("unexpected")
			}
			return err
		}); err != nil {
			t.Fatalf("'%s': %v", buf, err)
		}
		if nGot != n {
			t.Fatalf("'%s': %d (got) != %d (expected)",
				buf, nGot, n,
			)
		}
		if sGot != str {
			t.Fatalf("'%s': %q (got) != %q (expected)",
				buf, sGot, str,
			)
		}
	})
}
