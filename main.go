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
var updatedMap = make(map[string]string)

func removeFirstAlphanumeric(s string) string {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+`)
	return re.ReplaceAllString(s, "")
}

func main() {
x2:
	for key, value := range channelstwo {
		newKey := strings.Replace(key, "_", " ", -1)
		updatedMap[newKey] = value
	}
	fmt.Println(updatedMap)
	//renamer()
	convertList()
	for key, value := range updatedMap {
		channels[key] = value
	}
	//fmt.Println(channels)
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
	items = append(items, item("History"))
	items = append(items, item("Add a channel"))

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
			//match = strings.Replace(match, "-", "/", 1)
			// dobra jak to robie to cos sie dziwnego dzieje z tytulami z jakiegos powodu
			// aaa moge chyba zwyczajnie w tym drugim zamienic myslniki na jakis inny znak i bedzie git
			match = removeFirstAlphanumeric(match)
			match = strings.TrimSpace(match)
			//fmt.Println(match)
			mecze = append(mecze, match)
			// mecze to jest lista tytulow jak cos

			videos[match] = link
			video = video + 1
		}
	} else {
		// notatka ze przeciez mam jeszcze ta mape co przy robieniu history sie zrobila
		// dobra tu mam mapke filmikow wiec pewnie trzeba zrobic zeby liczylo ile ma elementow i na tej podstawie dostosowac dlugosc w opisie frameworka
		// also ciekawe bo jak dam na poczatku jakis kanal a potem historie to sie wyswietla normalnie historia tylko potem wszystko inne to tez historia
		// !!! dobra teraz jak tak patrze na te historie to to w ogole nie jest to co powinno byc w sensie output.txt bo nie ma niektorych filmikow
		// a nie jednak jest tylko w jakis dziwny sposob to dziala bo np jeden filmik pokazuje sie dopiero na 10 stronie
		// ok czyli samo titleLinkMap jest w porzadku tylko to sie laduje jakos w nieskonczonosc chyba bo jest do 312
		fmt.Println(itemshist)
		// !!! dobra czyli generalnie to jest problem z tym ze itemshist sie zapetla w ktoryms momencie chyba
		videos = titleLinkMap

		fmt.Println(len(itemshist))
		//listHeight = len(videos)
		// a moze trzeba czesc tego kodu w ogole przeniesc do maina zamiast trzymac we frameworku
	}

	for i, str := range mecze {
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
		mecze[i] = strings.ReplaceAll(str, "&#34;", "")
	}

	for i := range mecze {
		// z tym bledem na dole to cos tutaj w sumie moze byc albo w sumie nie bo to goto drugim w mainie i tak mnie cofa do poczatku wiec moze trzeba jakos zrobic zeby przywracalo wartosci framework.go do tych co byly na poczatku
		itemstwo = append(itemstwo, item(mecze[i]))
		// tylko mecze to jest w ogole ta tablica wiec w teorii ten kod jest w ogole niepotrzebny
		// itemki jest uzywane 3 razy i potem nic sie z tym nie dzieje
		itemki = append(itemki, mecze[i])
	}
	// ! tutaj moge cos pokombinowac zeby w history bylo rozwiazane jakos identycznie
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
	// sprobuje jutro wyklikac ze wiecej filmikow bo moze one dlatego sie zapetlaja
	/*
		dobra czyli ogolnie kolejnosc jest zalatwiona tylko teraz trzeba ogarnac to zapetlanie i cala funkcja historii bedzie zrobiona
	*/

	if isgb {
		// tutaj wyswietla filmiki tego kanalu co sie wybralo na poczatku nie wiem do konca czemu

		//z link = "" i redeklarowaniem modelu jest dalej to samo
		// jak cos to to nie jest problem z wartoscia linku tylko z lista prawdopodobnie
		// also w ogole teraz jak to naprawilem to play next video ekran sie nie odpala tylko wraca do wyboru filmiku a jak za pierwszym razem odpalam bez cofania to normalnie jest
		link = ""
		goto x2
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
	os.Remove("output.txt")
}
