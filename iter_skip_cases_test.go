package jir

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_skip(t *testing.T) {
	type testCase struct {
		ptr    interface{}
		inputs []string
	}
	var testCases []testCase

	testCases = append(testCases, testCase{
		ptr: (*string)(nil),
		inputs: []string{
			`""`,       // valid
			`"hello"`,  // valid
			`"`,        // invalid
			`"\"`,      // invalid
			`"\x00"`,   // invalid
			"\"\x00\"", // invalid
			"\"\t\"",   // invalid
			`"\t"`,     // valid
		},
	})
	testCases = append(testCases, testCase{
		ptr: (*[]interface{})(nil),
		inputs: []string{
			`[]`,             // valid
			`[1]`,            // valid
			`[  1, "hello"]`, // valid
			`[abc]`,          // invalid
			`[`,              // invalid
			`[[]`,            // invalid
		},
	})
	testCases = append(testCases, testCase{
		ptr: (*float64)(nil),
		inputs: []string{
			"+1",    // invalid
			"-a",    // invalid
			"-\x00", // invalid, zero byte
			"0.1",   // valid
			"0..1",  // invalid, more dot
			"1e+1",  // valid
			"1+1",   // invalid
			"1E1",   // valid, e or E
			"1ee1",  // invalid
			"100a",  // invalid
			"10.",   // invalid
		},
	})
	testCases = append(testCases, testCase{
		ptr: (*struct{})(nil),
		inputs: []string{
			`{}`,                         // valid
			`{"hello":"world"}`,          // valid
			`{hello:"world"}`,            // invalid
			`{"hello:"world"}`,           // invalid
			`{"hello","world"}`,          // invalid
			`{"hello":{}`,                // invalid
			`{"hello":{}}`,               // valid
			`{"hello":{}}}`,              // invalid
			`{"hello":  {  "hello": 1}}`, // valid
			`{abc}`,                      // invalid
		},
	})
	for _, testCase := range testCases {
		valType := reflect.TypeOf(testCase.ptr).Elem()
		for _, input := range testCase.inputs {
			t.Run(input, func(t *testing.T) {
				should := require.New(t)
				ptrVal := reflect.New(valType)
				stdErr := json.Unmarshal([]byte(input), ptrVal.Interface())
				iter := ParseString(Default, input)
				iter.Skip()
				iter.ReadNil() // trigger looking forward
				err := iter.Error
				if err == io.EOF {
					err = nil
				} else {
					err = errors.New("remaining bytes")
				}
				if stdErr == nil {
					should.Nil(err)
				} else {
					should.NotNil(err)
				}
			})
		}
	}
}
