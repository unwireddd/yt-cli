package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	//"os"

	"github.com/charmbracelet/bubbles/cursor"
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
var itemkihist []string
var itemshist []list.Item
var isHistory bool
var text2 string
var toSubs string
var nazwa string
var nazwaLink string
var linkingError string
var isgb bool
var lenHistory int
var linecounter int

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle2          = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type modelsix struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel3() modelsix {
	m := modelsix{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Invidious link"
			t.CharLimit = 64

		}

		m.inputs[i] = t
	}

	return m
}

func (m modelsix) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelsix) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *modelsix) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m modelsix) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	nazwa = m.inputs[0].Value()
	nazwaLink = m.inputs[1].Value()
	return b.String()
}

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
	ti.Placeholder = "[Invidious link]"
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
			isgb = true

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

	// to jest do zapisywania historii
	for key, value := range videos {
		if value == link {
			combinated := fmt.Sprintf("%s [Line break here] %s\n", key, link)
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

		if _, err := tea.NewProgram(initialModel3()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}

		if !strings.HasPrefix(nazwaLink, "https://") {
			fmt.Println("This doesnt seem to be a valid invidious link")
			os.Exit(1)
		}

		channels[nazwa] = nazwaLink

		nazwa = strings.ReplaceAll(nazwa, " ", "_")
		toSubs = fmt.Sprintf("%s	%s ", nazwa, nazwaLink)
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
		link = channels[nazwa]
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
		// output.txt to history z usunietymi duplikatami jak cos
		file, _ := os.Open("output.txt")
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, " [Line break here] ", 2)
			if len(parts) != 2 {

				continue
			}

			title := parts[0]
			link := parts[1]

			titleLinkMap[title] = link

			itemkihist = append(itemkihist, parts[0])

		}

		for i, j := 0, len(itemkihist)-1; i < j; i, j = i+1, j-1 {
			itemkihist[i], itemkihist[j] = itemkihist[j], itemkihist[i]
		}

		countLines()
		fmt.Println(linecounter)
		itemkihist = itemkihist[:linecounter]

		for i := range itemkihist {
			// z tym bledem na dole to cos tutaj w sumie moze byc albo w sumie nie bo to goto drugim w mainie i tak mnie cofa do poczatku wiec moze trzeba jakos zrobic zeby przywracalo wartosci framework.go do tych co byly na poczatku
			itemshist = append(itemshist, item(itemkihist[i]))
			itemki = append(itemki, itemkihist[i])
		}
		itemshist = itemshist[:linecounter]

		lenHistory = len(titleLinkMap)

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
				combinated := fmt.Sprintf("%s [Line break here] %s\n", key, link)
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
				combinated := fmt.Sprintf("%s [Line break here] %s\n", key, link)
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
