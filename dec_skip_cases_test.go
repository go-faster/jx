package jx

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_skip(t *testing.T) {
	type testCase struct {
		ptr    interface{}
		name   string
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
			"0",      // valid
			"+1",     // invalid
			"-a",     // invalid
			"-\x00",  // invalid, zero byte
			"0.1",    // valid
			"0e1",    // valid
			"0e+1",   // valid
			"0e-1",   // valid
			"0e",     // invalid
			"-e",     // invalid
			"+e",     // invalid
			".e",     // invalid
			"0.e",    // invalid
			"0.e",    // invalid
			"0.0e",   // invalid
			"0.0e1",  // valid
			"0.0e+1", // valid
			"0.0e+",  // invalid
			"0.0e-",  // invalid
			"0..1",   // invalid, more dot
			"1e+1",   // valid
			"1+1",    // invalid
			"1E1",    // valid, e or E
			"1ee1",   // invalid
			"100a",   // invalid
			"10.",    // invalid
			"-0.12",  // valid
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
		t.Run(valType.Kind().String(), func(t *testing.T) {
			for inputIdx, input := range testCase.inputs {
				t.Run(fmt.Sprintf("Test%d", inputIdx), func(t *testing.T) {
					t.Cleanup(func() {
						if t.Failed() {
							t.Logf("Input: %q", input)
						}
					})
					should := require.New(t)
					ptrVal := reflect.New(valType)
					stdErr := json.Unmarshal([]byte(input), ptrVal.Interface())
					iter := DecodeStr(input)
					if stdErr == nil {
						should.NoError(iter.Skip())
						should.ErrorIs(iter.Null(), io.ErrUnexpectedEOF)
					} else {
						should.Error(func() error {
							if err := iter.Skip(); err != nil {
								return err
							}
							if err := iter.Skip(); err != io.EOF {
								return err
							}
							return nil
						}())
					}
				})
			}
		})
	}
}
