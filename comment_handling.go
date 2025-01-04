package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/anaskhan96/soup"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type modelc struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modelc) Init() tea.Cmd {
	return nil
}

func (m modelc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelc) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

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

	items := []list.Item{}

	for _, match := range matches {
		items = append(items, item(match[1]))
	}

	l := list.New(items, itemDelegate{}, 20, listHeight)
	l.Title = "Comments"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := modelc{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// dobra teraz plan jest taki zeby te wszystkie komentarze dac w liste zeby sie ladnie wyswietlalo
}
