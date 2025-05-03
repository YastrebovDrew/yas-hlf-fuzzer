package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------
// флаги
// -----------------------------------------------------------------------------
var (
	manifestPath = flag.String("manifest", "", "path to manifest.json")
	outDir       = flag.String("out", "corpus", "directory for seed files")
	limitTotal   = flag.Int("limit", 0, "max total files (0 = unlimited)")
	limitPerFunc = flag.Int("limit-per-func", 0, "max files per function (0 = unlimited)")
	format       = flag.String("format", "pipe", "seed format: pipe|json")
)

// -----------------------------------------------------------------------------
// структура манифеста
// -----------------------------------------------------------------------------
type Manifest struct {
	Functions []struct {
		Name string     `json:"name"`
		Args [][]string `json:"args"`
	} `json:"functions"`
}

// -----------------------------------------------------------------------------
// main
// -----------------------------------------------------------------------------
func main() {
	flag.Parse()
	if *manifestPath == "" {
		fail("flag -manifest is required")
	}
	raw, err := os.ReadFile(*manifestPath)
	check(err)

	var mf Manifest
	check(json.Unmarshal(raw, &mf))

	check(os.MkdirAll(*outDir, 0o755))

	total := 0
	index := 0
	for _, fn := range mf.Functions {
		count := 0
		emit := func(payload []string) {
			// respect per-func лимит
			if *limitPerFunc > 0 && count >= *limitPerFunc {
				return
			}
			// respect глобальный лимит
			if *limitTotal > 0 && total >= *limitTotal {
				return
			}
			file := filepath.Join(*outDir, fmt.Sprintf("%06d_%s", index, fn.Name))
			writeSeed(file, payload)
			index++
			count++
			total++
		}
		product(fn.Name, fn.Args, emit)
	}
	fmt.Printf("generated %d files in %s (%s format)\n", total, *outDir, *format)
}

// -----------------------------------------------------------------------------
// рекурсивно строим декартово произведение
// -----------------------------------------------------------------------------
func product(name string, matrix [][]string, emit func([]string)) {
	var recur func(int, []string)
	recur = func(i int, acc []string) {
		if i == len(matrix) {
			emit(append([]string{name}, acc...))
			return
		}
		for _, v := range matrix[i] {
			recur(i+1, append(acc, v))
		}
	}
	recur(0, nil)
}

// -----------------------------------------------------------------------------
// запись сида
// -----------------------------------------------------------------------------
func writeSeed(path string, parts []string) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	switch *format {
	case "json":
		b, _ := json.Marshal(parts)
		f.Write(b)

	case "pipe":
		for i, p := range parts {
			if i > 0 {
				io.WriteString(f, "|")
			}
			io.WriteString(f, p)
		}

	default:
		fail("unknown -format; use pipe or json")
	}
}

// -----------------------------------------------------------------------------
// util
// -----------------------------------------------------------------------------
func check(err error) {
	if err != nil {
		fail(err.Error())
	}
}
func fail(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}
