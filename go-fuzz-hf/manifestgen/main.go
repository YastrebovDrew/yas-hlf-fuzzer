package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Manifest описывает контракт: список функций и их аргументыs.
type Manifest struct {
	Functions []Function `json:"functions"`
}

type Function struct {
	Name string     `json:"name"`
	Args [][]string `json:"args"`
}

func main() {
	out := flag.String("out", "manifest.json", "Путь к manifest.json или директории для него")
	flag.Parse()

	// Расширяем '~' в начале пути
	path := *out
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[1:])
		}
	}

	// Если out указывает на директорию, создаём файл manifest.json в ней
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		path = filepath.Join(path, "manifest.json")
	}

	// Убеждаемся, что директория для файла существует
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания каталога %s: %v\n", dir, err)
		os.Exit(1)
	}

	// Пример манифеста
	example := Manifest{
		Functions: []Function{
			{Name: "InitLedger", Args: [][]string{}},
			{Name: "CreateAsset", Args: [][]string{
				{"asset1", "asset_😀", "<long_id>", ""},
				{"blue", "red", ""},
				{"0", "1", "-1"},
				{"alice", "bob", ""},
				{"0", "1000", "-999"},
			}},
			{Name: "ReadAsset", Args: [][]string{{"asset1", "unknown", ""}}},
		},
	}

	data, err := json.MarshalIndent(example, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка маршалинга манифеста: %v\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("Шаблон manifest.json создан в %s\n", path)
}
