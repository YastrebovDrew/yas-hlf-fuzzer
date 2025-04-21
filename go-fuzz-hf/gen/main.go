// ddd
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Manifest struct {
	Functions []Function `json:"functions"`
}

type Function struct {
	Name string     `json:"name"`
	Args [][]string `json:"args"`
}

func main() {
	manifestPath := flag.String("manifest", "manifest.json", "Путь к manifest.json")
	outDir := flag.String("out", "corpus", "Каталог для сидов")
	limit := flag.Int("limit", 200, "Максимум файлов на функцию")
	flag.Parse()

	raw, err := ioutil.ReadFile(*manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения %s: %v\n", *manifestPath, err)
		os.Exit(1)
	}
	var m Manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка разбора %s: %v\n", *manifestPath, err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания %s: %v\n", *outDir, err)
		os.Exit(1)
	}

	idx := 0
	for _, fn := range m.Functions {
		combos := generateCombos(fn.Name, fn.Args, *limit)
		for _, combo := range combos {
			fname := fmt.Sprintf("%04d_%s", idx, fn.Name)
			path := filepath.Join(*outDir, fname)
			if err := ioutil.WriteFile(path, []byte(combo), 0644); err != nil {
				panic(fmt.Sprintf("Не могу записать %s: %v", path, err))
			}
			idx++
		}
	}
	fmt.Printf("Сгенерировано %d сидов в '%s'\n", idx, *outDir)
}

func generateCombos(fnName string, args [][]string, limit int) []string {
	if len(args) == 0 {
		return []string{fnName}
	}
	var res []string
	current := make([]string, len(args))
	count := 0
	var dfs func(int)
	dfs = func(pos int) {
		if count >= limit {
			return
		}
		if pos == len(args) {
			combo := fnName
			for _, v := range current {
				combo += "\x00" + v
			}
			res = append(res, combo)
			count++
			return
		}
		for _, val := range args[pos] {
			if count >= limit {
				break
			}
			current[pos] = val
			dfs(pos + 1)
		}
	}
	dfs(0)
	return res
}
