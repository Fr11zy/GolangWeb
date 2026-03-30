package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	//"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output io.Writer, path string, flag bool) error {
	return printTree(output, path, "", flag)
}

func printTree(output io.Writer, path, prefix string, flag bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := file.ReadDir(-1)
	if err != nil {
		return err
	}

	if !flag {
		n := 0
		for _, entry := range entries {
			if entry.IsDir() {
				entries[n] = entry
				n++
			}
		}
		entries = entries[:n]
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for ind, entry := range entries {
		isLast := (ind == len(entries)-1)
		var symbol string
		if isLast {
			symbol = "└───"
		} else {
			symbol = "├───"
		}

		if entry.IsDir() {
			fmt.Fprintf(output, "%s%s%s\n", prefix, symbol, entry.Name())
		}
		if !entry.IsDir() && flag {
			info, _ := entry.Info()
			size := info.Size()
			if size == 0 {
				fmt.Fprintf(output, "%s%s%s (empty)\n", prefix, symbol, entry.Name())
			} else {
				fmt.Fprintf(output, "%s%s%s (%db)\n", prefix, symbol, entry.Name(), size)
			}
		}

		if entry.IsDir() {
			var newprefix string
			if isLast {
				newprefix = prefix + "\t"
			} else {
				newprefix = prefix + "│\t"
			}

			printTree(output, filepath.Join(path, entry.Name()), newprefix, flag)
		}
	}
	return nil
}
