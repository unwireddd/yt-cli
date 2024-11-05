package main

import (
	"bufio"
	"fmt"
	"os"
)

// jak se wstawilem ta funkcje do framework.go to tez chyba nie dzialalo

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
	fmt.Sprintln("The line is", lineCount)

	linecounter = lineCount
}
