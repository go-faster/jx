package misc_tests

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	j "github.com/ogen-go/json"
)

func Test_read_big_float(t *testing.T) {
	should := require.New(t)
	iter := j.ParseString(j.ConfigDefault, `12.3`)
	val := iter.ReadBigFloat()
	val64, _ := val.Float64()
	should.Equal(12.3, val64)
}

func Test_read_big_int(t *testing.T) {
	should := require.New(t)
	iter := j.ParseString(j.ConfigDefault, `92233720368547758079223372036854775807`)
	val := iter.ReadBigInt()
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())
}

func Test_read_number(t *testing.T) {
	should := require.New(t)
	iter := j.ParseString(j.ConfigDefault, `92233720368547758079223372036854775807`)
	val := iter.ReadNumber()
	should.Equal(`92233720368547758079223372036854775807`, string(val))
}

func Test_encode_inf(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.Inf(1))
	should.Error(err)
	_, err = json.Marshal(float32(math.Inf(1)))
	should.Error(err)
	_, err = json.Marshal(math.Inf(-1))
	should.Error(err)
}

func Test_encode_nan(t *testing.T) {
	should := require.New(t)
	_, err := json.Marshal(math.NaN())
	should.Error(err)
	_, err = json.Marshal(float32(math.NaN()))
	should.Error(err)
	_, err = json.Marshal(math.NaN())
	should.Error(err)
}

func Benchmark_json_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := float64(0)
		_ = json.Unmarshal([]byte(`1.1`), &result)
	}
}
