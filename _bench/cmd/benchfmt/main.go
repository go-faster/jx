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
	type row struct {
		group    string
		subgroup string
		name     string
		cols     []string
	}
	var rows []row

	for scanner.Scan() {
		line := scanner.Text()
		m := benchRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Remove "Benchmark" prefix from name
		name := strings.TrimPrefix(m[1], "Benchmark")
		parts := strings.SplitN(name, "/", 4)
		group, subgroup, benchname := "", "", ""
		if len(parts) > 1 {
			group = parts[0]
			if len(parts) > 2 {
				subgroup = parts[1]
				benchname = strings.Join(parts[2:], "/")
			} else {
				benchname = parts[1]
			}
		} else {
			group = name
		}
		rows = append(rows, row{
			group:    group,
			subgroup: subgroup,
			name:     benchname,
			cols:     []string{name, m[3], m[4], m[5], m[6]},
		})
	}

	lastGroup, lastSub := "", ""
	for _, r := range rows {
		if r.group != lastGroup {
			fmt.Printf("\n## %s\n\n", r.group)
			lastGroup = r.group
			lastSub = ""
		}
		if r.subgroup != "" && r.subgroup != lastSub {
			fmt.Printf("### %s\n\n", r.subgroup)
			// Always print table header for each subtable
			fmt.Println("| Benchmark | ns/op | MB/s | B/op | allocs/op |")
			fmt.Println("|-----------|-------|------|------|-----------|")
			lastSub = r.subgroup
		}
		if r.subgroup == "" && (lastSub != "" || r.name != "") {
			// Print table header for top-level group without subgroups
			fmt.Println("| Benchmark | ns/op | MB/s | B/op | allocs/op |")
			fmt.Println("|-----------|-------|------|------|-----------|")
			lastSub = ""
		}
		fmt.Printf("| %s | %s | %s | %s | %s |\n",
			r.name, r.cols[1], r.cols[2], r.cols[3], r.cols[4])
	}
}
