//go:build go1.8
// +build go1.8

package misc_tests

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	j "github.com/ogen-go/json"
)

func Test_read_uint64_invalid(t *testing.T) {
	should := require.New(t)
	iter := j.ParseString(j.ConfigDefault, ",")
	iter.ReadUint64()
	should.NotNil(iter.Error)
}

func Test_float_as_int(t *testing.T) {
	should := require.New(t)
	var i int
	should.NotNil(json.Unmarshal([]byte(`1.1`), &i))
}

func Benchmark_jsoniter_encode_int(b *testing.B) {
	stream := j.NewStream(j.ConfigDefault, ioutil.Discard, 64)
	for n := 0; n < b.N; n++ {
		stream.Reset(nil)
		stream.WriteUint64(0xffffffff)
	}
}

func Benchmark_itoa(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strconv.FormatInt(0xffffffff, 10)
	}
}

func Benchmark_jsoniter_int(b *testing.B) {
	iter := j.NewIterator(j.ConfigDefault)
	input := []byte(`100`)
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		iter.ReadInt64()
	}
}

func Benchmark_json_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := int64(0)
		_ = json.Unmarshal([]byte(`-100`), &result)
	}
}
