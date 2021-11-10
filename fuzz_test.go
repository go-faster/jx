//go:build go1.18
// +build go1.18

package jx

import (
	"bytes"
	"testing"

	"github.com/go-faster/errors"
)

func FuzzValid(f *testing.F) {
	for _, s := range []string{
		"{}",
		`{"foo": "bar"}`,
		``,
		`"foo"`,
		`"{"`,
		`"{}"`,
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		Valid(data)
	})
}

func FuzzDecEnc(f *testing.F) {
	f.Add([]byte("{}"))
	f.Add([]byte(`"foo"`))
	f.Add([]byte(`123"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"foo": {"bar": 1, "baz": [1, 2, 3]}}`))
	f.Add([]byte(`"\u003cf\ufffd\ufffd\ufffd"`))
	f.Add([]byte(`"a\ufffdz"`))
	f.Add([]byte(`"\\nH\\tel\\tl\\ro\\\\World\\r\n\rHello\r\tHi"`))
	f.Add([]byte(`"key:\"/registry/runtimeclasses/\" range_end:\"/registry/runtimeclasses0\" count_only:true "`))
	f.Fuzz(func(t *testing.T, data []byte) {
		r := GetDecoder()
		r.ResetBytes(data)
		defer PutDecoder(r)

		v, err := r.Any()
		if err != nil {
			t.Skip()
		}
		if v.Type == AnyInvalid {
			t.Skip()
		}
		w := GetEncoder()
		w.Any(v)

		// Parsing from buf to new value.
		r.ResetBytes(w.Bytes())
		parsed, err := r.Any()
		if err != nil {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s\nData:  %s",
				err, w.Bytes(), v, data)
		}
		if !parsed.Equal(v) {
			t.Fatalf("%v:\nBuf:   %s\nValue: %s != %s \nData:  %s",
				nil, w.Bytes(), parsed, v, data)
		}
		b := w.Bytes()
		w.SetBytes(nil)
		parsed.Write(w)
		if !bytes.Equal(w.Bytes(), b) {
			t.Fatalf("%s != %s", w, b)
		}
	})
}

func FuzzValues(f *testing.F) {
	f.Add(int64(1), "hello")
	f.Add(int64(1534564316421), " привет ")
	f.Fuzz(func(t *testing.T, n int64, str string) {
		w := GetEncoder()
		defer PutEncoder(w)

		w.ArrStart()
		w.Int64(n)
		w.Str(str)
		w.ArrEnd()

		i := GetDecoder()
		i.ResetBytes(w.Bytes())
		var (
			nGot int64
			sGot string
		)
		if err := i.Arr(func(i *Decoder) error {
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
			t.Fatalf("'%s': %v", w, err)
		}
		if nGot != n {
			t.Fatalf("'%s': %d (got) != %d (expected)",
				w, nGot, n,
			)
		}
		if sGot != str {
			t.Fatalf("'%s': %q (got) != %q (expected)",
				w, sGot, str,
			)
		}
	})
}
