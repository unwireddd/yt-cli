package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/anaskhan96/soup"
)

func loadComments(url string) {
	//url := "https://inv.nadeko.net/watch?v=Oii2zmEqWVU&nojs=1"

	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(resp)

	comments := doc.Find("div", "id", "comments").HTML()

	re := regexp.MustCompile(`<p\s+style="white-space:pre-wrap">(.*?)</p>`)

	// Find all matches
	matches := re.FindAllStringSubmatch(comments, -1)

	//re := regexp.MustCompile(`<p>(.*?)</p>`)

	//matches := re.FindAllString(comments, -1)

	for _, match := range matches {
		fmt.Println(match[1])
		fmt.Println("")
	}
}
