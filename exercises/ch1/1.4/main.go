package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	lineToFiles := make(map[string]map[string]int)

	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup1.4: %v\n", err)
			continue
		}

		input := bufio.NewScanner(f)
		for input.Scan() {
			text := input.Text()
			counts[text]++

			if _, ok := lineToFiles[text]; !ok {
				lineToFiles[text] = make(map[string]int)
			}

			lineToFiles[text][filename] = 1
		}

		f.Close()
	}

	for line, count := range counts {
		if count > 1 {
			var files string
			for file := range lineToFiles[line] {
				files += file + " "
			}

			fmt.Printf("%d\t%s\t%s\n", count, files, line)
		}
	}
}
