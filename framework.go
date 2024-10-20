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
	"time"

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

// end test modelfour

// test modelfive

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

	// We handle errors just like any other message
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

// end test modelfive

func restart() {
	args := os.Args
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second) // wait for the new process to start
	os.Exit(0)
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
		// tutaj sie dodaje text w sensie nazwa kanalu i jakis losowy link niezalezny od niczego
		// mozna w sumie zrobic osobny parser dla kanalow i pokazuje pierwszy znaleziony
		// albo w sumie chuj w ostatecznosci moge tez pytac o link i w ten sposob dodawac i zapisywac to w tekscie a potem appendowac za kadzym odpaleniem
		toSubs = fmt.Sprintf("%s	%s ", text, text2)
		// poki co jest w takim formacie ale to nic mozna potem cos wykombinowac ai zapytac zeby to zmienic na format mapkowy
		// teraz najwiekszy problem to te jebane zapytania we frameworku co sie wyswietlaja 2 razy z jakiegos powodu
		// to zapytanie o kanal w ogole chyba 2 razy to zapisuje w pliku jak sie cos faktycznie wpisze 2 razy w zapytaniu
		file, err := os.OpenFile("channels.md", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Write the content to the file
		_, err = file.WriteString(toSubs)
		//_, err = file.WriteString("\n")

		//_, err = file.WriteString("\n")
		// zeby to dobrze dzialalo to trzeba zeby kazdy zestaw link + kanal byl razem w osobnej linii i wtedy bedzie git

		// !!! mozna zmienic format pliku z .txt na .md to wtedy bedzie czytac htmlowy syntax tylko nie wiadomo czy skrypt bedzie to obslugiwal wtedy
		// i chuj niby dalej jest to w nowej linii nie wiem czemu

		// !!!
		// moge teoretycznie napisac osobny skrypt do przerabiania pliku .md zeby byl wlasciwy format ale nie wiem czy to dobry pomysl
		// dobra to dziala i jest normalnie w output.md

		// dalej jest zle ale mozna sie spytac ai zeby mi
		// text2 to link ktory faktycznie sie gdzies zapisuje dopiero w drugim zapytaniu
		// text to nazwa ktora zapisuje sie w pierwszym zapytaniu
		if err != nil {
			panic(err)
		}
		// dobra i teraz to mozna jakos zapisac do jakiegos pliku i napisac zeby ladowalo to do channels.go przy kazdym odpaleniu programu
		link = channels[text]
		// dziala jak cos teraz tylko trzeba wykombinowac jak to dodac do mapy w channels.go na stale
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

		/*
			l := list.New(itemshist, itemDelegate{}, 17, listHeight)
			l.Title = "Select the channel you'd like to watch"
			l.SetShowStatusBar(false)
			l.SetFilteringEnabled(false)
			l.Styles.Title = titleStyle
			l.Styles.PaginationStyle = paginationStyle
			l.Styles.HelpStyle = helpStyle

			m = model{list: l} // tutaj jestem w modelu model i on sam siebie wywoluje

			if _, err := tea.NewProgram(m).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		*/

		for title, link := range titleLinkMap {
			fmt.Printf("%s: %s\n", title, link)
		}
		link = "https://iv.nboeck.de/channel/UC7YOGHUfC1Tb6E4pudI9STA"

		// to nie bedzie dzialac bo w mainie jest przeciez scrapowanie tego wszystkiego wiec trzeba pewnie zrobic jakis if statement i tam mecze to beda inaczej wyciagane w ogole

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
				//fmt.Printf("The key associated with the value '%s' is '%s'\n", link, key)
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
