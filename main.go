package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
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
var itemsgb []list.Item
var updatedMap = make(map[string]string)
var sprawdzam []string
var howgb = 0
var sprawdzanieczegos int
var mapConvert []string
var isQuittin bool
var items []list.Item
var isAppending int
var checksForGoingBack bool
var howManyAdded = 61
var howManyAddedS = 20

// cel na najblizsze dnie to zrobienie parsowania dla nastepnych stron przy szukaniu

// filmiki sie nie wysietlaja znowu z jakiegos powodu

//var userAgents map[string]string

func removeFirstAlphanumeric(s string) string {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+`)
	return re.ReplaceAllString(s, "")
}

func main() {
	//userAgents["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	//Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	isVideoLoading = false
	//loading:
	if isVideoLoading {
		fmt.Println("You are now in the full video loading mode, loading the full video list may take a while")
	}

	// notatka ze nastawienie na poczatku isgb na false nie dziala ale chyba trzeba cos wlasnie probowac w ta strone
	howgb = 0
	isAppending = 0

x2:

	isHistory = false

	var m tea.Model
	itemsthree := []list.Item{
		item("Replay video"),
		item("Play next video"),
		item("Play previous video"),
		item("Go back to videos list"),
		item("Display comments"),
	}
	var mecze []string

	if isAppending == 0 {
		for key, value := range channelstwo {
			newKey := strings.Replace(key, "_", " ", -1)
			updatedMap[newKey] = value
		}
		convertList()
		for key, value := range updatedMap {
			channels[key] = value
		}

		screen.Clear()

		for t := range maps.Keys(channels) {
			items = append(items, item(t))
		}

		items = append(items, item("Search"))
		items = append(items, item("History"))
		items = append(items, item("Add a channel"))
		items = append(items, item("Remove a channel"))

		isAppending++

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

	soup.Header("User-Agent", "User/Agent")
	ator, _ := soup.Get(link)
	atorvid := soup.HTMLParse(ator)
	linkTesting := link

loading:

	itemstwo = nil

	if isQuittin == true {
		return
	} else {

		if !isHistory {
			mecze = nil

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
				match = strings.Replace(match, `</p></a>`, "", -1)
				match = strings.Replace(match, `"`, "", -1)
				match = strings.Replace(match, `/`, fmt.Sprintf("https://www.youtube.com/"), -1)
				cutter := regexp.MustCompile(`https?://[^\s]+`)
				link := cutter.FindString(match)
				match = strings.Replace(match, `https://www.youtube.com/`, fmt.Sprintf(""), -1)
				match = strings.Replace(match, "watch?v=", "", 1)

				testlista := regexp.MustCompile(`\?list=.*?\s`)
				match = testlista.ReplaceAllString(match, "[Playlist]")

				if strings.Contains(match, "[Playlist]") {

					match = strings.ReplaceAll(match, "[Playlist]", "")
					match = fmt.Sprintf("%s [Playlist]", match)

				}

				match = strings.TrimPrefix(match, "?list=")
				match = removeFirstAlphanumeric(match)
				match = strings.TrimSpace(match)
				mecze = append(mecze, match)

				videos[match] = link
				video = video + 1

			}

			mecze = append(mecze, "Load all videos")

			if isVideoLoading && len(mecze) == 61 {

				howManyAdded += 10

				dlugosc := 61

				for len(mecze) == dlugosc && len(mecze) < howManyAdded {

					fmt.Println("Fetching videos")

					filmikidwa := atorvid.Find("div", "class", "page-next-container")

					kont := filmikidwa.Find("a")
					kont2 := kont.Attrs()
					sprawdzanie := kont2["href"]

					link = fmt.Sprint("https://inv.nadeko.net", sprawdzanie)
					ator, _ = soup.Get(link)
					atorvid = soup.HTMLParse(ator)
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
						match = strings.Replace(match, `</p></a>`, "", -1)
						match = strings.Replace(match, `"`, "", -1)
						match = strings.Replace(match, `/`, fmt.Sprintf("https://www.youtube.com/"), -1)
						cutter := regexp.MustCompile(`https?://[^\s]+`)
						link := cutter.FindString(match)
						match = strings.Replace(match, `https://www.youtube.com/`, fmt.Sprintf(""), -1)
						match = strings.Replace(match, "watch?v=", "", 1)

						testlista := regexp.MustCompile(`\?list=.*?\s`)
						match = testlista.ReplaceAllString(match, "[Playlist]")

						if strings.Contains(match, "[Playlist]") {

							match = strings.ReplaceAll(match, "[Playlist]", "")
							match = fmt.Sprintf("%s [Playlist]", match)

						}

						match = strings.TrimPrefix(match, "?list=")
						match = removeFirstAlphanumeric(match)
						match = strings.TrimSpace(match)
						mecze = append(mecze, match)
						videos[match] = link
						video = video + 1
						dlugosc += 1

					}
				}
				mecze = mecze[len(mecze)-60:]
				mecze = append(mecze, "Load all videos")
				isVideoLoading = false

			} else if isVideoLoading && strings.Contains(linkTesting, "search") {

				// tutaj normalnie wykrywa to jak wejde w searcha

				strona := 2

				howManyAddedS += 10

				dlugosc := 20

				for len(mecze) == dlugosc && len(mecze) < howManyAddedS {

					fmt.Println("Fetching videos")

					//filmikidwa := atorvid.Find("div", "class", "page-next-container")

					//kont := filmikidwa.Find("a")
					//kont2 := kont.Attrs()
					//sprawdzanie := kont2["href"]

					page := fmt.Sprint("&page=", strona)

					link = fmt.Sprint(linkTesting, page)
					ator, _ = soup.Get(link)
					atorvid = soup.HTMLParse(ator)
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
						match = strings.Replace(match, `</p></a>`, "", -1)
						match = strings.Replace(match, `"`, "", -1)
						match = strings.Replace(match, `/`, fmt.Sprintf("https://www.youtube.com/"), -1)
						cutter := regexp.MustCompile(`https?://[^\s]+`)
						link := cutter.FindString(match)
						match = strings.Replace(match, `https://www.youtube.com/`, fmt.Sprintf(""), -1)
						match = strings.Replace(match, "watch?v=", "", 1)

						testlista := regexp.MustCompile(`\?list=.*?\s`)
						match = testlista.ReplaceAllString(match, "[Playlist]")

						if strings.Contains(match, "[Playlist]") {

							match = strings.ReplaceAll(match, "[Playlist]", "")
							match = fmt.Sprintf("%s [Playlist]", match)

						}

						match = strings.TrimPrefix(match, "?list=")
						match = removeFirstAlphanumeric(match)
						match = strings.TrimSpace(match)
						mecze = append(mecze, match)
						videos[match] = link
						video = video + 1
						dlugosc += 1
						strona += 1

					}
				}
				mecze = mecze[len(mecze)-20:]
				mecze = append(mecze, "Load all videos")
				isVideoLoading = false

			}

		} else {

			videos = titleLinkMap

		}

		if isgb {
			lendeleting := len(mecze)
			howgb += len(mecze)

			sprawdzam = mecze[len(mecze)-lendeleting:]

			itemsgb = itemsgb[:0]

		}

		for i, str := range mecze {
			mecze[i] = strings.ReplaceAll(str, "&#39;", "")
			mecze[i] = strings.ReplaceAll(str, "&#39;", "")
			mecze[i] = strings.ReplaceAll(str, "&#34;", "")
		}

		for i := range mecze {
			if isgb {

				itemsgb = append(itemsgb, item(sprawdzam[i]))

			}
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

			m = &modeltwo{list: l}

			if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

		} else {

			if isgb {

				l = list.New(itemsgb, itemDelegate{}, defaultWidth, listHeight)
				l.Title = "Select the video you'd like to watch"
				l.SetShowStatusBar(false)
				l.SetFilteringEnabled(false)
				l.Styles.Title = titleStyle
				l.Styles.PaginationStyle = paginationStyle
				l.Styles.HelpStyle = helpStyle

				m = &modeltwo{list: l}

				if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
			} else {
				if strings.Contains(link, "search") {

					sprawdzanieczegos = len(mecze)

					itemstwo = itemstwo[howgb:]
					finalremove := len(itemstwo) - len(mecze)
					itemstwo = itemstwo[finalremove:]

					if howgb > 0 {

					}

					itemki = itemki[howgb:]
					itemki = itemki[finalremove:]

				}

				l = list.New(itemstwo, itemDelegate{}, defaultWidth, listHeight)
				l.Title = "Select the video you'd like to watch"
				l.SetShowStatusBar(false)
				l.SetFilteringEnabled(false)
				l.Styles.Title = titleStyle
				l.Styles.PaginationStyle = paginationStyle
				l.Styles.HelpStyle = helpStyle

				m = &modeltwo{list: l}

				if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
				if isVideoLoading {
					goto loading
				}
			}
		}

	x3:
		if isReplaying == true {

			testt := exec.Command("mpv", linkForReplays)
			testt.Run()

		}

		if isgb {

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

		m = &modelthree{list: l}

		if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		if testowanie == "Go back to videos list" {
			isgb = true
			checksForGoingBack = true

			goto x
		}
		if testowanie == "Play next video" {
			goto x2
		}
		if testowanie == "Replay video" {
			goto x3
		}
		os.Remove("output.txt")
	}
}
