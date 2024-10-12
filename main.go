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

var lista []list.Item
var itemki []string
var itemstwo []list.Item

func removeFirstAlphanumeric(s string) string {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+`)
	return re.ReplaceAllString(s, "")
}

func main() {

	screen.Clear()
	var m tea.Model
	var items []list.Item
	itemsthree := []list.Item{
		item("Play next video"),
		item("Play previous video"),
		item("Go back to videos list"),
	}
	var mecze []string

	for t := range maps.Keys(channels) {
		items = append(items, item(t))
	}

	items = append(items, item("Search"))
	items = append(items, item("Test"))
	items = append(items, item("Add"))
	items = append(items, item("Add / Remove a channel"))
	items = append(items, item("History"))

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

	// tutaj trzeba zrobic if statement ze jak isHistory jest na false to sie robi to a jak nie to inaczej
	if !isHistory {
		ator, _ := soup.Get(link)
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
			cutter := regexp.MustCompile(`https?://[^\s]+`)
			link := cutter.FindString(match)
			match = strings.Replace(match, `https://www.youtube.com/`, fmt.Sprintf(""), -1)
			match = strings.Replace(match, "watch?v=", "", 1)
			match = strings.Replace(match, "?list=", "[Playlist]", 1)
			match = removeFirstAlphanumeric(match)
			match = strings.TrimSpace(match)

			mecze = append(mecze, match)

			videos[match] = link
			video = video + 1
		}
	} else {
		// notatka ze przeciez mam jeszcze ta mape co przy robieniu history sie zrobila
		videos = titleLinkMap
	}

	for i, str := range mecze {
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
		mecze[i] = strings.ReplaceAll(str, "&#34;", "")
	}

	for i := range mecze {
		itemstwo = append(itemstwo, item(mecze[i]))
		itemki = append(itemki, mecze[i])
	}
x:
	if isHistory {
		l = list.New(itemshist, itemDelegate{}, defaultWidth, listHeight)
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
	} else {
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

	l = list.New(itemsthree, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select option"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	screen.Clear()

	m = modelthree{list: l}

	if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if testowanie == "Go back to videos list" {
		goto x
	}

}
