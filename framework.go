// jutro te pointery z maila ogarnac

// JUTRO PIERWSZE CO ZROBIC TO OGARNAC COMMITOWANIE NA GH
// generalnie usuwanie kanalow wyglada na naprawione wiec teraz kolejnym priorytetem jest zrobienie wiekszej liczby filmikow ktore beda sie ladowac

// ten projekt jak cos to jest Projekty/yt-cli

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

// generalnie z tym przeskakiwaniem rzeczy na gorze w historii mozna zrobic cos takiego ze on za kazdym razem bedzie przy dodawaniu czegos skanowal cale output.txt czy jest tam cos o takiej samej nazwie i to
// jak chce cos innego robic akurat to git bo to remove a channel nie kloci sie z innymi czesciami kodu za bardzo generalnie

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
var testowaniedrugiejrzeczy string
var isReplaying bool
var linkForReplays string
var globalRmTest string
var isVideoLoading bool

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

func initialModel3() *modelsix {
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

	return &m
}

func (m *modelsix) Init() tea.Cmd {
	return textinput.Blink
}

func (m *modelsix) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

// tu jak jest pointer reciever to sie cos psuje tez jak cos

func (m *modelsix) View() string {
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

// tu sie zaczyna modelfour jak cos

type modelfour struct {
	textInput textinput.Model
	err       error
}

func initialModel() *modelfour {
	ti := textinput.New()
	// czyli widze ze w przyp
	ti.Placeholder = "Linux"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &modelfour{
		textInput: ti,
		err:       nil,
	}
}

func (m *modelfour) Init() tea.Cmd {
	return textinput.Blink
}

// czyli tutaj tez trzeba pokombinowac z pointerami zeby to wyswietlanie searcha 2 razy sie naprawilo i dodawania ale to jest troche inaczej skonstruowane niz tamte poprzednie z listami
// tez widze ze przez to ze tutaj jest initialmodel i przez niego to jest wywolywane
func (m *modelfour) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		// tutaj jest to do wychodzenia jak cos

		case tea.KeyCtrlC:

			// wazna rzecz ze tea.Quit sprawia ze program idzie dalej zamiast go wylaczac

			// zrobienie tutaj tego tak samo jak w tamtych innych z isgb true nie dziala jak cos pewnie przez strukture tego modelu ze inna
			//testowaniedrugiejrzeczy = "Yes"
			isQuittin = true
			isgb = true
			//dalej nie dziala ale tutaj to wyglada jakby w ogole tego nie wykrywalo z jakiegos powodu
			// trzeba zrobic cos takiego zeby tutaj sie wracalo albo zeby w jakis sposob wywalalo caly program jak to dam

			// ciekawe bo teraz dziala ale jest ten sam problem z wywalaniem calego terminala co wczesniej

			// os.Exit teoretycznie to mozna zrobic ale wtedy jest ten blad z wywalaniem sie terminala calego
			// w sumie jak nic nie dziala to mozna wykombinowac po prostu tak zeby sie wracalo kilkukrotnie to nie dziala przynajmniej jak daje tea.Quit kilka razy tutaj na gorze
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *modelfour) View() string {
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

func (m *modelfive) Init() tea.Cmd {
	return textinput.Blink
}

func (m *modelfive) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *modelfive) View() string {

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

func (m *modeltwo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	isgb = false
	// o nice tutaj dziala

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "q":
			isgb = true
			//m.choice = ""
			// wyzerowanie zmiennej link tez nie dziala
			// goto testowanko nie zadziala i wtedy w ogole sie wszystko zacina
			// przez defer tez nie bedzie dzialac bo nie moze wejsc do funkcji po zamknieciu jej
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				globaltest = m.choice
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *modeltwo) View() string {
	// tutaj jest pierwszy raz ten globaltest ktory sie przypisuje do filmiku z listy filmikow

	var ok bool
	filePath := "history"

	//fmt.Println("Test starts here")
	//fmt.Println(m.choice)
	//fmt.Println("Test ends here")

	// a moze to to updatea trzeba przeniesc czy cos
	// z testow wynika ze m.choice tutaj w ogole jeszcze nie mam tez ma to w sumie sens bo przeciez m.choice to jest to co wybralem w kanalach

	if m.choice == "Load all videos" {
		isVideoLoading = true

	} else {

		link, ok = videos[m.choice]

		//fmt.Println(videos[m.choice])
		//globaltest = m.choice
		/*fmt.Println("The globaltest variable")
		//fmt.Println(globaltest)
		fmt.Println(m.choice)
		fmt.Println("tests end")*/

		m.choice = ""

		// to jest do zapisywania historii

		// dobra czyli tutaj zamysl jest taki zeby przed tym jak to doda do historii to przeskanowac to cale i wywalic wszystkie duplikaty i wtedy powinno dzialac
		for key, value := range videos {
			if value == link {
				combinated := fmt.Sprintf("%s [Line break here] %s\n", key, link)
				historyCleanup(combinated)
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

			linkForReplays = link

			testt := exec.Command("mpv", link)
			testt.Run()
			return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
		}
		if m.quitting {
			return quitTextStyle.Render("Don't want to watch? That’s cool.")
		}
		// isgb = false
		// dobra czyli zrobienie tutaj isgb na false tez nie bedzie dzialalo bo tu jak jest q to nastawia isgb na true a potem na koncu sie i tak z tego robi false
	}

	return "\n" + m.list.View()

}

// START test mowy model do usuwania kanalow

type modelrm struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modelrm) Init() tea.Cmd {
	return nil
}

func (m modelrm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":

			isQuittin = true

			//m.quitting = true
			//return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				globalRmTest = m.choice
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelrm) View() string {
	var ok bool

	//fmt.Println(channelstwo["qwert"])
	// czyli wychodzi na to ze ten link w channelstwo w sensie mapce tutaj sie normalnie zgadza a usuwane dziala wiec to raczej jest problem z samym skryptem
	// also mozna chyba zwyczajnie usunac tez link wywolujac to drugi raz

	// dobra czyli tutaj pierwszy problem jest taki ze m.choice w calym modelu jest puste a drugi te cala konstrukcja tego modelu jest jakas dziwna i mysle ze cos tu sie nie zgadza

	// tutaj koniec tych wszystkich rzeczy ktore sa do wywalenia

	// dobra ogolnie wychodzi na to ze zeby to dzialalo to po pierwsze potrzebne jest to zeby to m.choice gdzies zapisywac a po drugie zeby poza nazwa usuwalo tez link

	//channelToRemove := fmt.Sprintln("%s    %s", "qwert", channelstwo["qwert"])

	//fmt.Println(channelToRemove)

	// oo dobra czyli teraz juz normalnie to dziala

	channelRemove(globalRmTest)
	channelRemove(channelstwo[globalRmTest])

	// Check for any errors

	if ok {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Don't want to watch? That’s cool.")
	}
	return "\n" + m.list.View()

}

// KONIEC test nowy model do usuwania kanalow

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
		case "q":

			isQuittin = true

			//m.quitting = true
			//return m, tea.Quit
		case "ctrl+c":
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

	if m.choice == "Add a channel" {

		testowanko()
		// z tym w ogole jest taki problem ze z tego co patrzylem po errorze to jak klikam w model add a channel to on z jakiegos powodu zaczyna od razu parsowanie zamiast wyswietlic ten model

		// czyli wychodzi na to ze to sie z jakiegos powodu nie odpala
		// podejrzewam ze problem jest w mainie a nie tutaj bo tutaj to nic nie zmienialem od dawna i sie dopiero po zrobieniu tego isgb zaczelo robic

		// tutaj mi executuje ten pierwszy model
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
		link = nazwaLink
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	// START test usuwanie kanalow

	if m.choice == "Remove a channel" {

		// na obecna chwile to nie dziala w sensie nie usuwa nic

		//fmt.Println(items)

		// dobra czyli sam zamysl na usuwanie kanalow z tej listy jest taki ze chce wywolac liste jak w pierwszej tylko ten kanal co jest wybrany usunac z channels.md chyba
		// also generalnie to mysle ze do usuwania kanalow mozna sprobowac zrobic osobny model ktory bedzie dzialal jak model zwykly tylko bez tych wszystkich rzeczy z opcjami a z samymi kanalami

		const defaultWidth = 20

		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = "Select the channel you'd like to remove"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		// tutaj to widze ze znowu jest ten sam problem co juz milion razy byl

		// ciekawe w sumie bo widze ze tutaj moge se zmienic m na inna nazwe i normalnie dziala bez tego errora

		// jest niby opcja ze to jest problem z wywolywaniem modelu w innym modelu ale z drugiej strony to jest normalnie ten sam problem co juz kiedys byl wiec pewnie nie
		// teraz jest kolejny problem z tym ze ta lista jest pusta ale to nie wiem czy bardziej wina tego co zrobilem w sensie przypisaniu m jako nowa lokalna zmienna czy tego ze nic nie zrobilem z dodawaniem tej listy
		// items to jest niby globalna zmienna ale jak ja na gorze printuje to jest pusta

		m := modelrm{list: l}

		// ten problem tutaj jest najprawdopodobniej przez to ze w pliku history sie robi jakis whitespace nie wiem czemu
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	// KONIEC test usuwanie kanalow

	if m.choice == "Search" {
		isgb = false
		isHistory = false
		// jak tutaj nastawie isgb na false to jest znowu ten problem z tym ze na liscie najpierw sie wyswietlaja filmiki z poprzednich a potem dopiero szukana fraza a jak nie to jest out of range
		// tutaj raczej sie nie da nic wykombinowac bo to tylko nastawia link a cale parsowanie jest w mainie
		screen.Clear()

		// teraz z tym znowu jest jakis problem ze jak wejde w jakis kanal a potem z niego do searcha to sie wyswietlaja filmiki z niego co juz kiedys chyba bylo naprawiane tylko nie pamietam w jaki sposob
		// no widze ze w tej nowej wersji tez jest ten problem nie wiem w ogole jak to sie stalo ze wszystkie bugfixy poprzednie nagle jakos zniknely

		// gwiazdka to jest pointer reciever btw
		// dobra ale widze ze przy szukaniu te pointery nic w ogole nie daja
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		if testowaniedrugiejrzeczy == "Yes" {
			log.SetOutput(io.Discard)
			log.Panic("Quitting...")
		}

		encodedText := url.QueryEscape(text)
		link = "https://inv.nadeko.net/search?q=" + encodedText
		//fmt.Println(link)
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "History" {

		// tutaj tez teoretycznie os.Exit powinno rozwiazac sprawe ale bedzie ten sam problem co z innymi rzeczami

		isgb = false
		// zamienianie isgb tutaj chyba tez nic nie daje jak cos ale idk

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
		//itemkihist = itemkihist[:len(itemkihist)-1]

		countLines()
		//fmt.Println(linecounter)
		itemkihist = itemkihist[:linecounter]

		for i := range itemkihist {
			// z tym bledem na dole to cos tutaj w sumie moze byc albo w sumie nie bo to goto drugim w mainie i tak mnie cofa do poczatku wiec moze trzeba jakos zrobic zeby przywracalo wartosci framework.go do tych co byly na poczatku
			itemshist = append(itemshist, item(itemkihist[i]))
			itemki = append(itemki, itemkihist[i])
		}
		itemshist = itemshist[:linecounter]

		lenHistory = len(titleLinkMap)

		link = "https://inv.nadeko.net/channel/UC7YOGHUfC1Tb6E4pudI9STA"

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

func (m *modelthree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			// notatka ze to tutaj jest ten problem z tym ze caly terminal sie wywala
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
			// o zmienienie jednego case na dwa osobne chyba zadzialalo
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

func (m *modelthree) View() string {

	isReplaying = false

	if m.choice == "Play next video" {

		// dobra czyli teraz to teoretycznie dziala ale jest problem z listami ze sie te same filmiki pokazuja

		//02.12.2024
		// ciekawe bo jak sie wracam to sie zwieksza liczba kanalow w pierwszej liscie
		// ale nvm bo to dziala generalnie
		// zeby to naprawic to trzeba wywalic to appendowanie w kolko jak cos
		isgb = true
		howgb += 60
		testowanie = m.choice
		// generalnie w tym calym to problem jest taki ze wartosc m.choice czyli tym samym globaltest jest pusta bo cos tam sie psuje na gorze przy view tej funkcji

		// START to cale jest w ogole do zignorowania bo tylko przypisuje rzeczy do historii

		for key, value := range videos {

			// dobra czyli tutaj zamysl jest taki zeby przed tym jak to doda do historii to przeskanowac to cale i wywalic wszystkie duplikaty i wtedy powinno dzialac
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

		// KONIEC

		var index int
		// jak cos to itemki to jest normalna tablica a nie lista wiec z jakims dziwnym przeskakiwaniem w kolejnosci tez nie powinno byc problemow
		// also widze ze jak jednak szukam w samym frontendzie to tez te filmiki jakos dziwnie przeskakuja wiec tu problem moze byc w tym ze on to jakos jeszcze raz skanuje przy play next czy cos

		// wychodzi na to ze globaltest jest puste co by wyjasnialo dlaczego nie dziala
		// prawdopodobnie to wina tego pointera bo wczesniej to normalnie dzialalo jak powinno
		for i, value := range itemki {
			if value == globaltest {

				// czyli to nastepne to zwyczajnie zawsze jest 1 niewazne co wiec z odwracaniem tablicy opcja odpada

				// globaltest to juz jest w tym drugim m.choice z modeltwo w ktorym jest link w sensie normalne szukanie czyli z tym jest niby wszystko git
				index = i
				break
			}
		}

		//fmt.Println(index)
		// w tym momencie index jest 0 czyli tak jak powinno byc a po dodaniu jest 1 czyli tez teoretycznie dobrze
		// a jak juz dam filmik 6 to tez jest 0 i 1 z jakiegos powodu
		index = index + 1

		// zawsze jest na 0 i 1 i z tego co widze to teraz np sie w ogole pokazal ten sam filmik
		//fmt.Println(index)
		m.choice = itemki[index]
		link = videos[m.choice]
		testt := exec.Command("mpv", link)
		testt.Run()

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "Play previous video" {
		// START historia do zignorowania
		isgb = true
		testowanie = m.choice

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

		// KONIEC

		var index int
		for i, value := range itemki {
			if value == globaltest {
				index = i
				break
			}
		}

		if index == 0 {
			// dobra widze czyli tu jest ten problem co u gory ze jest zawsze 0 i 1 jak printuje to z jakiegos powodu
			return quitTextStyle.Render(fmt.Sprintf("%s? I can't see any previous video.", m.choice))
		}
		//fmt.Println(index)
		// w tym momencie index jest 0 czyli tak jak powinno byc a po dodaniu jest 1 czyli tez teoretycznie dobrze
		// a jak juz dam filmik 6 to tez jest 0 i 1 z jakiegos powodu
		index = index - 1
		//fmt.Println(index)
		m.choice = itemki[index]
		link = videos[m.choice]
		testt := exec.Command("mpv", link)
		testt.Run()
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	if m.choice == "Go back to videos list" {
		// ciekawe w sumie bo przy go back to sie w ogole gdzies zapetla i caly czas jest
		// jak bede ogarniac zapetlanie tego to mozna popatrzec na to potem
		testowanie = m.choice
		return "ok"
	}

	if m.choice == "Replay video" {
		// ciekawe w sumie bo przy go back to sie w ogole gdzies zapetla i caly czas jest
		// jak bede ogarniac zapetlanie tego to mozna popatrzec na to potem
		isReplaying = true
		testowanie = m.choice
		return "ok"
	}
	if m.quitting {
		return quitTextStyle.Render("Don't want to watch? That’s cool.")
	}
	return "\n" + m.list.View()
}
