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

var testBools = []string{
	"",
	"tru",
	"fals",
	"fal\x00e",
	"fals\x00",
	"f\x00\x00\x00\x00",
	"nope",
	"true",
	"false",
	" true",
	" false",
	"\ntrue",
	"\nfalse",
	"\ttrue",
	"\tfalse",
}

var testNumbers = append([]string{
	"",                            // invalid
	"0",                           // valid
	"-",                           // invalid
	"--",                          // invalid
	"+",                           // invalid
	".",                           // invalid
	"e",                           // invalid
	"E",                           // invalid
	"-.",                          // invalid
	"-1",                          // valid
	"--1",                         // invalid
	"+1",                          // invalid
	"++1",                         // invalid
	"-a",                          // invalid
	"-0",                          // valid
	"00",                          // invalid
	"01",                          // invalid
	".00",                         // invalid
	"00.1",                        // invalid
	"-00",                         // invalid
	"-01",                         // invalid
	"-\x00",                       // invalid, zero byte
	"0.1",                         // valid
	"0e0",                         // valid
	"-0e0",                        // valid
	"+0e0",                        // valid
	"0e-0",                        // valid
	"-0e-0",                       // valid
	"+0e-0",                       // valid
	"0e+0",                        // valid
	"-0e+0",                       // valid
	"+0e+0",                       // valid
	"0e+01234567890123456789",     // valid
	"0.00e-01234567890123456789",  // valid
	"-0e+01234567890123456789",    // valid
	"-0.00e-01234567890123456789", // valid
	"0e1",                         // valid
	"0e+1",                        // valid
	"0e-1",                        // valid
	"0e-11",                       // valid
	"0e-1a",                       // invalid
	"1.e1",                        // invalid
	"0e-1+",                       // invalid
	"0e",                          // invalid
	"0.e",                         // invalid
	"0-e",                         // invalid
	"0e-",                         // invalid
	"0e+",                         // invalid
	"0.0e",                        // invalid
	"0.0e1",                       // valid
	"0.0e+",                       // invalid
	"0.0e-",                       // invalid
	"0e0+0",                       // invalid
	"0.e0+0",                      // invalid
	"0.0e+0",                      // valid
	"0.0e+1",                      // valid
	"0.0e0+0",                     // invalid
	"e",                           // invalid
	"-e",                          // invalid
	"+e",                          // invalid
	".e",                          // invalid
	"e.",                          // invalid
	"0.",                          // invalid
	"1.",                          // invalid
	"0..1",                        // invalid, more dot
	"0.1.",                        // invalid, more dot
	"1..",                         // invalid, more dot
	"1e+1",                        // valid
	"1+1",                         // invalid
	"1E1",                         // valid, e or E
	"1ee1",                        // invalid
	"100a",                        // invalid
	"10.",                         // invalid
	"-0.12",                       // valid
	"0]",                          // invalid
	"0e]",                         // invalid
	"0e+]",                        // invalid
	"1.2.3",                       // invalid
	"0.0.0",                       // invalid
	"9223372036854775807",         // valid
	"9223372036854775808",         // valid
	"9223372036854775807.1",       // valid
	" 9223372036854775807",        // valid
	" 9223372036854775808",        // valid
	" 9223372036854775807.1",      // valid
	"\n9223372036854775807",       // valid
	"\n9223372036854775808",       // valid
	"\n9223372036854775807.1",     // valid
}, []string{
	// Test cases from strconv.

	// Copyright 2009 The Go Authors. All rights reserved.
	// Use of this source code is governed by a BSD-style
	// license that can be found in the LICENSE file.
	"1e23",
	"1E23",
	"100000000000000000000000",
	"1e-100",
	"123456700",
	"99999999999999974834176",
	"100000000000000000000001",
	"100000000000000008388608",
	"100000000000000016777215",
	"100000000000000016777216",
	"1e-20",
	"625e-3",

	// zeros
	"0",
	"0e0",
	"-0e0",
	"+0e0",
	"0e-0",
	"-0e-0",
	"+0e-0",
	"0e+0",
	"-0e+0",
	"+0e+0",
	"0e+01234567890123456789",
	"0.00e-01234567890123456789",
	"-0e+01234567890123456789",
	"-0.00e-01234567890123456789",

	"0e291",
	"0e292",
	"0e347",
	"0e348",
	"-0e291",
	"-0e292",
	"-0e347",
	"-0e348",

	// NaNs
	"nan",
	"NaN",
	"NAN",

	// Infs
	"inf",
	"-Inf",
	"+INF",
	"-Infinity",
	"+INFINITY",
	"Infinity",

	// largest float64
	"1.7976931348623157e308",
	"-1.7976931348623157e308",

	// next float64 - too large
	// "1.7976931348623159e308",  // json.Valid checks overflow.
	// "-1.7976931348623159e308", // json.Valid checks overflow.

	// the border is ...158079
	// borderline - okay
	"1.7976931348623158e308",
	"-1.7976931348623158e308",
	// borderline - too large
	// "1.797693134862315808e308",  // json.Valid checks overflow.
	// "-1.797693134862315808e308", // json.Valid checks overflow.

	// a little too large
	"1e308",
	// "2e308", // json.Valid checks overflow.
	// "1e309", // json.Valid checks overflow.

	// way too large
	// "1e310",     // json.Valid check overflow.
	// "-1e310",    // json.Valid check overflow.
	// "1e400",     // json.Valid check overflow.
	// "-1e400",    // json.Valid check overflow.
	// "1e400000",  // json.Valid check overflow.
	// "-1e400000", // json.Valid check overflow.

	// denormalized
	"1e-305",
	"1e-306",
	"1e-307",
	"1e-308",
	"1e-309",
	"1e-310",
	"1e-322",
	// smallest denormal
	"5e-324",
	"4e-324",
	"3e-324",
	// too small
	"2e-324",
	// way too small
	"1e-350",
	"1e-400000",

	// try to overflow exponent
	"1e-4294967296",
	// "1e+4294967296", // json.Valid check overflow.
	"1e-18446744073709551616",
	// "1e+18446744073709551616", // json.Valid check overflow.

	// Parse errors
	"1e",
	"1e-",
	".e-1",
	"1\x00.2",

	// https://www.exploringbinary.com/java-hangs-when-converting-2-2250738585072012e-308/
	"2.2250738585072012e-308",
	// https://www.exploringbinary.com/php-hangs-on-numeric-value-2-2250738585072011e-308/
	"2.2250738585072011e-308",

	// A very large number (initially wrongly parsed by the fast algorithm).
	"4.630813248087435e+307",

	// A different kind of very large number.
	"22.222222222222222",
	"2." + strings.Repeat("2", 4000) + "e+1",

	// Exactly halfway between 1 and math.Nextafter(1, 2).
	// Round to even (down).
	"1.00000000000000011102230246251565404236316680908203125",
	// Slightly lower; still round down.
	"1.00000000000000011102230246251565404236316680908203124",
	// Slightly higher; round up.
	"1.00000000000000011102230246251565404236316680908203126",
	// Slightly higher, but you have to read all the way to the end.
	"1.00000000000000011102230246251565404236316680908203125" + strings.Repeat("0", 10000) + "1",

	// Halfway between x := math.Nextafter(1, 2) and math.Nextafter(x, 2)
	// Round to even (up).
	"1.00000000000000033306690738754696212708950042724609375",

	// Halfway between 1090544144181609278303144771584 and 1090544144181609419040633126912
	// (15497564393479157p+46, should round to even 15497564393479156p+46, issue 36657)
	"1090544144181609348671888949248",
	// slightly above, rounds up
	"1090544144181609348835077142190",
}...)

var testStrings = append([]string{
	`""`,                   // valid
	`"hello"`,              // valid
	`"`,                    // invalid
	`"foo`,                 // invalid
	`"\`,                   // invalid
	`"\"`,                  // invalid
	`"\u`,                  // invalid
	`"\u1`,                 // invalid
	`"\u12`,                // invalid
	`"\u123`,               // invalid
	`"\u\n"`,               // invalid
	`"\u1\n"`,              // invalid
	`"\u12\n"`,             // invalid
	`"\u12\n"`,             // invalid
	`"\u123\n"`,            // invalid
	`"\u1d`,                // invalid
	`"\u$`,                 // invalid
	`"\21412`,              // invalid
	`"\uD834\1`,            // invalid
	`"\uD834\u3`,           // invalid
	`"\uD834\`,             // invalid
	`"\uD834`,              // invalid
	`"\u07F9`,              // invalid
	`"\u1234\n"`,           // valid
	`"\x00"`,               // invalid
	"\"\x00\"",             // invalid
	"\"\t\"",               // invalid
	"\"\\b\x06\"",          // invalid
	`"\t"`,                 // valid
	`"\n"`,                 // valid
	`"\r"`,                 // valid
	`"\b"`,                 // valid
	`"\f"`,                 // valid
	`"\/"`,                 // valid
	`"\\"`,                 // valid
	"\"\\u000X\"",          // invalid
	"\"\\uxx0X\"",          // invalid
	"\"\\uxxxx\"",          // invalid
	"\"\\u000.\"",          // invalid
	"\"\\u0000\"",          // valid
	"\"\\ua123\"",          // valid
	"\"\\uffff\"",          // valid
	"\"\\ueeee\"",          // valid
	"\"\\uFFFF\"",          // valid
	`"ab\n` + "\x00" + `"`, // invalid
	`"\n0123456"`,
}, func() (r []string) {
	// Generate tests for invalid space sequences.
	for i := byte(0); i <= ' '; i++ {
		r = append(r, `"`+string(i)+`"`)
	}
	// Generate tests to ensure unroll correctness.
	for i := byte('0'); i <= '9'; i++ {
		// Generates "0", "11", "222" ...
		n := int(i - '0' + 1)
		str := strings.Repeat(string(i), n)
		r = append(r, `"`+str+`"`)
		if n > 2 {
			// Insert newline.
			// Generates "22\n2", "333\n3" ...
			str = str[:n-2] + `\n` + str[n-1:]
			r = append(r, `"`+str+`"`)
			// Generates "2\n22", "3\n333" ...
			str = str[:1] + `\n` + str[2:]
			r = append(r, `"`+str+`"`)
		}
	}
	return r
}()...)

var testObjs = []string{
	"",                              // invalid
	"nope",                          // invalid
	"nul",                           // invalid
	"fals e",                        // invalid
	"nil",                           // invalid
	`{`,                             // invalid
	`{}`,                            // valid
	`{"1}`,                          // invalid
	`{"1:}`,                         // invalid
	`{"1,}`,                         // invalid
	`{"1":}`,                        // invalid
	`{"\1":}`,                       // invalid
	`{"1",}`,                        // invalid
	`{"1":,}`,                       // invalid
	`{"hello":"world"}`,             // valid
	`{hello:"world"}`,               // invalid
	`{"hello:"world"}`,              // invalid
	`{"hello","world"}`,             // invalid
	`{"hello":{}`,                   // invalid
	`{"hello":{}}`,                  // valid
	`{"hello":{}}}`,                 // invalid
	`{"hello":  {  "hello": 1}}`,    // valid
	`{abc}`,                         // invalid
	`invalid`,                       // invalid
	`{"foo"`,                        // invalid
	`{"foo"bar`,                     // invalid
	`{"foo": "bar",`,                // invalid
	`{"foo": "bar", true`,           // invalid
	`{"foo": "bar", "bar":`,         // invalid
	`{"foo": "bar", "bar":t`,        // invalid
	`{"foo": "bar", "bar":true`,     // invalid
	`{"foo": "bar", "bar"false`,     // invalid
	`{"foo": "bar", "bar": "bar"""`, // invalid
	`{"foo":`,                       // invalid
	`{"foo": "bar"`,                 // invalid
	`{"foo": "bar`,                  // invalid
	`{"foo": "bar",}`,               // invalid
	`{"foo": "bar", true}`,          // invalid
	"{\n\"foo\"\n: \n10e1   \n, \n\"bar\"\n: \ntrue\n}", // valid
}

var testArrs = []string{
	`[]`,             // valid
	`[1]`,            // valid
	`[  1, "hello"]`, // valid
	`[abc]`,          // invalid
	`[`,              // invalid
	`[,`,             // invalid
	`[[]`,            // invalid
	"[true,f",        // invalid
	"[true",          // invalid
	"[true,",         // invalid
	"[true]",         // invalid
	"[true,]",        // invalid
	"[true,false",    // invalid
	"[true,false,",   // invalid
	"[true,false,]",  // invalid
	"[true,false}",   // invalid
}

func TestDecoder_Skip(t *testing.T) {
	type testCase struct {
		ptr    interface{}
		inputs []string
	}
	var testCases []testCase

	testCases = append(testCases, testCase{
		ptr:    (*bool)(nil),
		inputs: testBools,
	})
	testCases = append(testCases, testCase{
		ptr:    (*string)(nil),
		inputs: testStrings,
	})
	numberCase := testCase{
		ptr:    (*float64)(nil),
		inputs: testNumbers,
	}
	testCases = append(testCases, numberCase)
	arrayCase := testCase{
		ptr:    (*[]interface{})(nil),
		inputs: append([]string(nil), testArrs...),
	}
	for _, c := range numberCase.inputs {
		arrayCase.inputs = append(arrayCase.inputs, `[`+c+`]`)
	}
	testCases = append(testCases, arrayCase)
	testCases = append(testCases, testCase{
		ptr:    (*struct{})(nil),
		inputs: testObjs,
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
