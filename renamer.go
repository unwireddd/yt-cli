package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func testowanko() {
	inputFile, err := os.Open("channels.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create("output.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				fmt.Fprintf(outputFile, "%s    %s\n", fields[i], fields[i+1])
			} else {
				fmt.Fprintf(outputFile, "%s\n", fields[i])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
