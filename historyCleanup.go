package main

import (
	//"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func historyCleanup(video string) {
	filePath := "history"
	strToRemove := video

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var content []byte
	content, err = io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !strings.Contains(string(content), strToRemove) {
		return
	}

	newContent := strings.ReplaceAll(string(content), strToRemove, "")

	err = file.Truncate(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write([]byte(newContent))
	if err != nil {
		fmt.Println(err)
		return
	}
}
