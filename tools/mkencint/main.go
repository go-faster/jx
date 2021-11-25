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
	Name       string
	Iterations int
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

var (
	//go:embed gen.tmpl
	rawTemplate string
)

func computeIterations(max uint64) int {
	// Compute maximum pow of 1000 plus remainder.
	return len(strconv.FormatUint(max, 10))/3 + 1
}

func generate(w io.Writer, pkgName string) error {
	buf := bytes.Buffer{}

	types := []IntType{
		{
			Name:       "int64",
			Iterations: computeIterations(math.MaxUint64),
		},
		{
			Name:       "int32",
			Iterations: computeIterations(math.MaxUint32),
		},
		{
			Name:       "int16",
			Iterations: computeIterations(math.MaxUint16),
		},
	}

	t := template.Must(template.New("gen").Funcs(template.FuncMap{
		"times": times,
		"title": title,
		"add":   add,
		"sub":   sub,
	}).Parse(rawTemplate))
	if err := t.ExecuteTemplate(&buf, "main", Config{
		PackageName: pkgName,
		Types:       types,
	}); err != nil {
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

func run() error {
	var (
		o       = flag.String("output", "", "output file")
		pkgName = flag.String("package", "jx", "package name")
	)
	flag.Parse()

	var w io.Writer = os.Stdout
	if path := *o; path != "" {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			fmt.Println(f.Close())
		}()
		w = f
	}

	return generate(w, *pkgName)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
