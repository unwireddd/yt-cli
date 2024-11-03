package main

import (
	"bufio"
	"fmt"
	"os"
)

func countLines() {
	file, err := os.Open("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	linecounter = lineCount
}
