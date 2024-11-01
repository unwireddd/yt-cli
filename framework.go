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
var itemshist []list.Item
var isHistory bool
var text2 string
var toSubs string
var nazwa string
var nazwaLink string
var linkingError string
var isgb bool
var lenHistory int

// testing podwojne

/*
	-na to does not implement tea model bylo jakies latwe rozwiazanie bo to robilem juz wczesniej z poprzednimi modelami (pewnie cos w mainie tylko musze to znalezc)
	-albo porownac z tym drugim szukaniem bo oba nie sa wykonywane w mainie
*/

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

// w initialModel sie nie robi problem a w tym na dole juz tak
// also wiadomo juz ze to nie wina samego kodu w przykladzie bo jak testowalem to dzialalo normalnie
// nie wina helpStyle ani initialmodel3

/*
-dziala i to byla literowka jakas
-teraz ogolnie trzeba przerobic tak zeby to zapisywalo w zmienna i robilo to samo co robia tamte dwa ktore sa osobno
*/

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

		// Change cursor mode
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

		// Set focus to next input

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
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
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			//fmt.Println(m.inputs[0].Value())

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *modelsix) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
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

// START wpisywanie

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
	// if empty
	text2 = m.textInput.Value()
	return fmt.Sprintf(
		"Search:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

// !!! dobra tutaj jest jednak wysokosc listy

const listHeight = 30

// jak to ustawie na 15 to w listach dalej jest 30 nie wiem czemu

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
		// na dole jest ustawianie ile elementow jest na liscie chyba
		// we wszystkich przypadkach jest max 30 elementow a tu jest 100 bo w innych jest 30 filmikow a tu ciagnie z mojej listy z history
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			// to jest do listy filmikow jak cos chyba ze to w ogole nie jest z modeltwo tylko z tym pierwszym
			isgb = true
			//m.quitting = true
			//return m, tea.Quit

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
			//os.Exit(0)
			// teraz teoretycznie sie nie wraca tylko normalnie wychodzi ale jest ten blad z wywalaniem sie terminala
			// dobra nvm to jednak nie jest problem chyba
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

	/*
		-juz za pierwszym dodaniem tego przy szukaniu byl problem ze sie odpala 2 razy
		-mozna uzyc tego https://github.com/charmbracelet/bubbletea/blob/main/examples/textinputs/main.go

	*/

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

		/*p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		p = tea.NewProgram(initialModel2())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		*/

		channels[nazwa] = nazwaLink
		// jak nazwa kanalu ma wiecej niz jedno slowo to sie psuje ale to mozna jakos naprawic np cudzyslow
		// also ten skrypt do zamiany mi usuwa duplikaty z channels.md wiec nie trzeba sie tym martwic w ogole

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
		// jak tu dam history to jest w ogole to samo
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
			// tutaj mozna w sumie cos pokombinowac z tym zmienianiem dlugosci listy zamiast w mainie bo to framework i tu deklaruje modele wiec jednak latwiej
			// notatka ze to jest model pierwszy wiec mozna cos normalnie zapisac w variabla i zamienic w modeltwo a modeltwo jest na gorze
			// albo w sumie mozna cos poprobowac z
			title := parts[0]
			link := parts[1]

			titleLinkMap[title] = link

			itemshist = append(itemshist, item(parts[0]))

			// aa czyli moge teraz sobie tez tutaj zrobic taka liste tytulow i potem ja jakos matchowac z ta mapka i wtedy powinno byc po kolei
			// dobra czyli teraz jest git tylko sie wyswietla na odwrot wiec pewnie trzeba odwrocic tablice

			//itemshist = append(itemshist, item(titleLinkMap[title]))
			// jak tutaj se printuje titleLinkMap to sie pokazuje kilka razy wiec mozliwe ze to sie zapetla jakos
			//fmt.Println(len(titleLinkMap))
			// test petla
			// tutaj jak probuje zeby dodal tylko pierwsze 30 to dalej sie wypierdala

		}
		histlen := len(itemshist) / 2
		itemshist = itemshist[:histlen]

		// normalny len itemshist to 32 a jak tutaj zrobie dzielenie przez 2 to nagle len jest 12

		for i, j := 0, len(itemshist)-1; i < j; i, j = i+1, j-1 {
			itemshist[i], itemshist[j] = itemshist[j], itemshist[i]
		}

		// jak podziele przez 2 to w ogole sie wyswietla tylko peirwsze 11 zamiast 15 z jakiegos powodu also to chyba nie jest problem z dlugoscia tablicy tylko z tym jak on to laduje potem
		// nawet jak dlugosc tablicy to jest 16 to i tak wyswietla 18
		// i nie wiedziec czemu w history sie tez jakos losowo wyswietla kanal

		// i chyba w ogole zapetla sie dokladnie 2 razy
		// teraz chyba jest git tylko trzeba jeszcze jakos wykombinowac zeby sie przestalo zapetlac

		/*for t := range maps.Keys(titleLinkMap) {
			itemshist = append(itemshist, item(t))
		} */

		lenHistory = len(titleLinkMap)

		/*for title, link := range titleLinkMap {
			fmt.Printf("%s: %s\n", title, link)
		} */

		link = "https://iv.nboeck.de/channel/UC7YOGHUfC1Tb6E4pudI9STA"

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	// wyglada na to ze z tym kodem jest wszystko git tylko cos sie psuje w tym wyswietlaniu go

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
