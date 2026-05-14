package src

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"
)

type stats struct {
	lines int
	files int
}

type Counter struct {
	results                map[string]*stats
	totalLines, totalFiles int
}

func NewCounter() *Counter {
	c := Counter{
		results: make(map[string]*stats),
	}

	return &c
}

func (c *Counter) Start() {
	err := filepath.WalkDir(".", c.walkFunc)

	if err != nil {
		fmt.Printf("Ошибка при обходе директории: %v\n", err)
		return
	}

	c.printResults()
}

func (c *Counter) walkFunc(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		if slices.Contains(ignoreDirs, d.Name()) {
			return filepath.SkipDir
		}
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := extensions[ext]; ok {
		count, err := c.countLines(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения %s: %v\n", path, err)
			return nil
		}

		if _, exists := c.results[lang]; !exists {
			c.results[lang] = &stats{}
		}

		c.results[lang].lines += count
		c.results[lang].files++
		c.totalLines += count
		c.totalFiles++
	}

	return nil
}

func (c *Counter) countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

func (c *Counter) printResults() {
	if c.totalFiles == 0 {
		fmt.Println("Файлы поддерживаемых языков не найдены.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Language\tFiles\tLines")
	fmt.Fprintln(w, "--------\t-----\t-----")

	// Сортировка по количеству строк (убывание)
	keys := make([]string, 0, len(c.results))
	for k := range c.results {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return c.results[keys[i]].lines > c.results[keys[j]].lines
	})

	for _, lang := range keys {
		fmt.Fprintf(w, "%s\t%d\t%d\n", lang, c.results[lang].files, c.results[lang].lines)
	}

	fmt.Fprintln(w, "--------\t-----\t-----")
	fmt.Fprintf(w, "TOTAL\t%d\t%d\n", c.totalFiles, c.totalLines)
	w.Flush()
}
