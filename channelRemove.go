// ten skrypt na dole z jakiegos powodu usunal mi cale output.md

// z tym to sie chyba w ogole cos innego dzieje bo mam tu the linux experiment ktory wczesniej sie usunal z outputa przez ten skrypt

// dobra chuj sprobuje jeszcze raz se to wygenerowac

package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func channelRemove(channel string) {
	// Define the string variable
	strToRemove := channel

	// Define the file path
	filePath := "output.md"

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the data to a string
	fileContent := string(data)

	// ogolnie to tutaj mozna normalnie do strToRemove dodac link ktory jest po nazwie kanalu

	// Remove all occurrences of the string
	newContent := strings.ReplaceAll(fileContent, strToRemove, "")

	// Write the new content to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Content removed successfully.")
}
