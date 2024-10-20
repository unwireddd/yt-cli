package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"maps"
	"net/url"
	"os"
	"os/exec"
	"strings"

	//"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inancgumus/screen"
)

var link string
var text string
var video int
var videos map[string]string = make(map[string]string)
var globaltest string
var testowanie string
var titleLinkMap = make(map[string]string)
var itemshist []list.Item
var isHistory bool
var text2 string
var toSubs string

// test model4

type (
	errMsg error
)

type modelfour struct {
	textInput textinput.Model
	err       error
}

func initialModel() modelfour {
	ti := textinput.New()
	ti.Placeholder = "Linux"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return modelfour{
		textInput: ti,
		err:       nil,
	}
}

func (m modelfour) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelfour) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m modelfour) View() string {
	text = m.textInput.Value()
	return fmt.Sprintf(
		"Search:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

type modelfive struct {
	textInput textinput.Model
	err       error
}

func initialModel2() modelfive {
	ti := textinput.New()
	ti.Placeholder = "Linux"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return modelfive{
		textInput: ti,
		err:       nil,
	}
}

func (m modelfive) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelfive) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m modelfive) View() string {
	if m.textInput.Value() == "" {
		return "This pool cannot be empty"
	}
	text2 = m.textInput.Value()
	return fmt.Sprintf(
		"Search:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

const listHeight = 30

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#C724B1"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type modeltwo struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modeltwo) Init() tea.Cmd {
	return nil
}

func (m modeltwo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m modeltwo) View() string {
	globaltest = m.choice
	var ok bool
	filePath := "history"

	link, ok = videos[m.choice]

	for key, value := range videos {
		if value == link {
			//fmt.Printf("The key associated with the value '%s' is '%s'\n", link, key)
			combinated := fmt.Sprintf("%s - %s\n", key, link)
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			_, err = io.WriteString(file, combinated)
			if err != nil {
				panic(err)
			}
		}
	}

	if ok {
		testt := exec.Command("mpv", link)
		testt.Run()
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Don't want to watch? That’s cool.")
	}
	return "\n" + m.list.View()
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
	var ok bool

	if m.choice == "Add / Remove a channel" {
		testowanko()
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		p = tea.NewProgram(initialModel2())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		channels[text] = text2
		toSubs = fmt.Sprintf("%s	%s ", text, text2)
		// for now the biggest problem is that those input models in framework are displayed 2 times Idk why
		file, err := os.OpenFile("channels.md", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(toSubs)

		if err != nil {
			panic(err)
		}
		link = channels[text]
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "Search" {
		screen.Clear()

		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		encodedText := url.QueryEscape(text)
		link = "https://iv.nboeck.de/search?q=" + encodedText
		fmt.Println(link)
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "History" {
		rmDuplicates()
		isHistory = true

		file, _ := os.Open("output.txt")
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) != 2 {
				//fmt.Printf("Skipping invalid line: %s\n", line)
				continue
			}

			title := parts[0]
			link := parts[1]

			titleLinkMap[title] = link
			for t := range maps.Keys(titleLinkMap) {
				itemshist = append(itemshist, item(t))
			}
		}

		for title, link := range titleLinkMap {
			fmt.Printf("%s: %s\n", title, link)
		}
		link = "https://iv.nboeck.de/channel/UC7YOGHUfC1Tb6E4pudI9STA"

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	link, ok = channels[m.choice]
	if ok {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Don't want to watch? That’s cool.")
	}
	return "\n" + m.list.View()
}

type modelthree struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modelthree) Init() tea.Cmd {
	return nil
}

func (m modelthree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			log.Fatal()
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

func (m modelthree) View() string {

	if m.choice == "Play next video" {

		for key, value := range videos {
			if value == link {
				combinated := fmt.Sprintf("%s - %s\n", key, link)
				file, err := os.OpenFile("history", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				_, err = io.WriteString(file, combinated)
				if err != nil {
					panic(err)
				}
			}
		}

		var index int
		for i, value := range itemki {
			if value == globaltest {
				index = i
				break
			}
		}
		index = index + 1
		m.choice = itemki[index]
		link = videos[m.choice]
		testt := exec.Command("mpv", link)
		testt.Run()
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "Play previous video" {

		for key, value := range videos {
			if value == link {
				combinated := fmt.Sprintf("%s - %s\n", key, link)
				file, err := os.OpenFile("history", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				_, err = io.WriteString(file, combinated)
				if err != nil {
					panic(err)
				}
			}
		}

		var index int
		for i, value := range itemki {
			if value == globaltest {
				index = i
				break
			}
		}

		if index == 0 {
			return quitTextStyle.Render(fmt.Sprintf("%s? I can't see any previous video.", m.choice))
		}
		index = index - 1
		m.choice = itemki[index]
		link = videos[m.choice]
		testt := exec.Command("mpv", link)
		testt.Run()
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "Go back to videos list" {
		testowanie = m.choice
		return "ok"
	}
	if m.quitting {
		return quitTextStyle.Render("Don't want to watch? That’s cool.")
	}
	return "\n" + m.list.View()
}
