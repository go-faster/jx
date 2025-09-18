package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var benchRe = regexp.MustCompile(`^(Benchmark\S+)\s+(\d+)\s+([\d.]+ ns/op)\s+([\d.]+ MB/s)\s+(\d+ B/op)\s+(\d+ allocs/op)`)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var rows [][]string

	for scanner.Scan() {
		line := scanner.Text()
		m := benchRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Remove "Benchmark" prefix from name
		name := strings.TrimPrefix(m[1], "Benchmark")
		rows = append(rows, []string{
			name, m[3], m[4], m[5], m[6],
		})
	}

	// Print markdown table without Iterations column.
	fmt.Println("| Benchmark | ns/op | MB/s | B/op | allocs/op |")
	fmt.Println("|-----------|-------|------|------|-----------|")
	for _, row := range rows {
		fmt.Printf("| %s | %s | %s | %s | %s |\n",
			row[0], row[1], row[2], row[3], row[4])
	}
}
