package main

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inancgumus/screen"
)

func main() {
	screen.Clear()
	var m tea.Model
	var items []list.Item
	var itemstwo []list.Item
	var mecze []string

	for t := range maps.Keys(channels) {
		items = append(items, item(t))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the channel you'd like to watch"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m = model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	ator, _ := soup.Get(url)
	atorvid := soup.HTMLParse(ator)
	owtput := atorvid.Find("div", "class", "pure-g").HTML()
	re := regexp.MustCompile(`<a[^>]*>.*?<p dir="auto">.*?</p>.*?</a>`)
	matches := re.FindAllString(owtput, -1)
	for _, match := range matches {

		match = strings.Replace(match, "<a href=", "", -1)
		if video < 10 {
			match = strings.Replace(match, `><p dir="auto">`, "			", -1)
		} else {
			match = strings.Replace(match, `><p dir="auto">`, "       ", -1)
		}
		//match = strings.Replace(match, `><p dir="auto">`, "			", -1)
		match = strings.Replace(match, `</p></a>`, "", -1)
		match = strings.Replace(match, `"`, "", -1)
		match = strings.Replace(match, `/`, fmt.Sprintf("https://www.youtube.com/"), -1)
		mecze = append(mecze, match)
		cutter := regexp.MustCompile(`https?://[^\s]+`)
		link := cutter.FindString(match)
		videos[match] = link
		video = video + 1
	}

	for i, str := range mecze {
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
	}

	for i := range mecze {
		itemstwo = append(itemstwo, item(mecze[i]))
	}

	l = list.New(itemstwo, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the video you'd like to watch"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m = modeltwo{list: l}

	if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
