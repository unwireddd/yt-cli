package main

import (
	"bufio"
	"fmt"
	"os"
)

func rmDuplicates() {
	// Open the file
	file, err := os.Open("history")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a map to store unique lines
	uniqueLines := make(map[string]bool)

	// Create a new file to write the unique lines
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line is already in the map
		if !uniqueLines[line] {
			uniqueLines[line] = true
			// Write the line to the output file
			_, err := outputFile.WriteString(line + "\n")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// Check for any errors that occurred during the scan
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
