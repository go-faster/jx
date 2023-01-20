// Command mkencint generates integer encoding/decoding functions.
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
	"strconv"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// IntType represents Go integer type.
type IntType struct {
	Name              string
	EncoderIterations int // ceil(log1000 (max value))
	DecoderIterations int // ceil(log10 (max value))
}

func defineIntType(name string, max uint64) IntType {
	formattedLen := len(strconv.FormatUint(max, 10))
	decoderIters := formattedLen

	const decoderItersLimit = 10 - 1
	if decoderIters > decoderItersLimit {
		decoderIters = decoderItersLimit
	}
	return IntType{
		Name:              name,
		EncoderIterations: formattedLen/3 + 1, // Compute maximum pow of 1000 plus remainder.
		DecoderIterations: decoderIters,       // Compute maximum pow of 10 plus remainder.
	}
}

var intTypes = []IntType{
	defineIntType("int8", math.MaxUint8),
	defineIntType("int16", math.MaxUint16),
	defineIntType("int32", math.MaxUint32),
	defineIntType("int64", math.MaxUint64),
}

// Config is generation config.
type Config struct {
	PackageName string
	Types       []IntType
}

func times(num int) []struct{} {
	return make([]struct{}, num)
}

func title(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func pow10(power int) (r int) {
	if power <= 0 {
		return 1
	}
	r = 10
	for i := 1; i < power; i++ {
		r *= 10
	}
	return r
}

func executeTemplate(w io.Writer, tmpl string, cfg Config) error {
	var buf bytes.Buffer

	t := template.Must(template.New("gen").Funcs(template.FuncMap{
		"times": times,
		"title": title,
		"add":   add,
		"sub":   sub,
		"pow10": pow10,
	}).Parse(tmpl))
	if err := t.ExecuteTemplate(&buf, "main", cfg); err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		_, _ = os.Stderr.Write(buf.Bytes())
		return fmt.Errorf("format: %w", err)
	}

	if _, err := w.Write(formatted); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

//go:embed encode.tmpl
var encodeTemplate string

func generateEncode(w io.Writer, pkgName string) error {
	return executeTemplate(w, encodeTemplate, Config{
		PackageName: pkgName,
		Types:       intTypes[1:], // Skip int8, use manual encoder.
	})
}

//go:embed decode.tmpl
var decodeTemplate string

func generateDecode(w io.Writer, pkgName string) error {
	return executeTemplate(w, decodeTemplate, Config{
		PackageName: pkgName,
		Types:       intTypes,
	})
}

func run() error {
	var (
		pkgName = flag.String("package", "jx", "package name")
	)
	flag.Parse()

	for _, file := range []struct {
		name string
		f    func(io.Writer, string) error
	}{
		{"w_int.gen.go", generateEncode},
		{"dec_int.gen.go", generateDecode},
	} {
		if err := func() error {
			f, err := os.Create(file.name)
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()

			return file.f(f, *pkgName)
		}(); err != nil {
			return fmt.Errorf("generate %s: %w", file.name, err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
