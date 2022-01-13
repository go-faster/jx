package jx

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestSkip(t *testing.T) {
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
	numberCase := testCase{
		ptr: (*float64)(nil),
		inputs: []string{
			"0",       // valid
			"-",       // invalid
			"+",       // invalid
			"-1",      // valid
			"+1",      // invalid
			"-a",      // invalid
			"-0",      // valid
			"-00",     // invalid
			"-01",     // invalid
			"-\x00",   // invalid, zero byte
			"0.1",     // valid
			"0e1",     // valid
			"0e+1",    // valid
			"0e-1",    // valid
			"0e-11",   // valid
			"0e-1a",   // invalid
			"1.e1",    // invalid
			"0e-1+",   // invalid
			"0e",      // invalid
			"e",       // invalid
			"-e",      // invalid
			"+e",      // invalid
			".e",      // invalid
			"e.",      // invalid
			"0.e",     // invalid
			"0-e",     // invalid
			"0e-",     // invalid
			"0e+",     // invalid
			"0.0e",    // invalid
			"0.0e1",   // valid
			"0.0e+",   // invalid
			"0.0e-",   // invalid
			"0e0+0",   // invalid
			"0.e0+0",  // invalid
			"0.0e+0",  // valid
			"0.0e+1",  // valid
			"0.0e0+0", // invalid
			"0..1",    // invalid, more dot
			"1e+1",    // valid
			"1+1",     // invalid
			"1E1",     // valid, e or E
			"1ee1",    // invalid
			"100a",    // invalid
			"10.",     // invalid
			"-0.12",   // valid
			"0]",      // invalid
			"0e]",     // invalid
			"0e+]",    // invalid
		},
	}
	testCases = append(testCases, numberCase)
	arrayCase := testCase{
		ptr: (*[]interface{})(nil),
		inputs: []string{
			`[]`,             // valid
			`[1]`,            // valid
			`[  1, "hello"]`, // valid
			`[abc]`,          // invalid
			`[`,              // invalid
			`[[]`,            // invalid
		},
	}
	for _, c := range numberCase.inputs {
		arrayCase.inputs = append(arrayCase.inputs, `[`+c+`]`)
	}
	testCases = append(testCases, arrayCase)
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

	testDecode := func(iter *Decoder, input string, stdErr error) func(t *testing.T) {
		return func(t *testing.T) {
			t.Cleanup(func() {
				if t.Failed() {
					t.Logf("Input: %q", input)
				}
			})

			should := require.New(t)
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
		}
	}
	for _, testCase := range testCases {
		valType := reflect.TypeOf(testCase.ptr).Elem()
		t.Run(valType.Kind().String(), func(t *testing.T) {
			for inputIdx, input := range testCase.inputs {
				t.Run(fmt.Sprintf("Test%d", inputIdx), func(t *testing.T) {
					ptrVal := reflect.New(valType)
					stdErr := json.Unmarshal([]byte(input), ptrVal.Interface())

					t.Run("Buffer", testDecode(DecodeStr(input), input, stdErr))

					r := strings.NewReader(input)
					d := Decode(r, 512)
					t.Run("Reader", testDecode(d, input, stdErr))

					r.Reset(input)
					obr := iotest.OneByteReader(r)
					t.Run("OneByteReader", testDecode(Decode(obr, 512), input, stdErr))
				})
			}
		})
	}
}
