// to normalnie dziala jak cos i zamienia to z channels.txt na poprawna mapke z go
/*
W sumie to do channels.txt mozna tez dodac subskrybcje z channels.go i jak mam plik tekstowy to niech on sie zwyczajnie konwertuje na plik go za kazdym odpaleniem
czy cos w tym stylu

also mozliwe ze on tez w sumie nie powinien mi generowac mapki z go tylko elementy i potem je jakos appendowac na samym poczatku go channels.go czy cos takiego
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func convertList() {
	file, err := os.Open("output.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var channels = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "https://", 2)
		if len(parts) != 2 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}
		key := strings.TrimSpace(parts[1])
		value := strings.TrimSpace(parts[0])
		channels[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	outputFile, err := os.Create("output.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString("package main")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = outputFile.WriteString("\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = outputFile.WriteString("var channelstwo = map[string]string{\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range channels {
		key = fmt.Sprintf("https://%s", key)
		_, err = outputFile.WriteString(fmt.Sprintf("    \"%s\":    \"%s\",\n", value, key))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	_, err = outputFile.WriteString("}\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}
