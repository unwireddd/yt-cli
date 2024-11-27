package main

// dobra czyli generalnie widze ze juz wszystko podstawowe co powinno dziala w kwestii historii i szukania wiec mozna wrocic na jakis czas do naprawiania bugow z poruszaniem sie po programie

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
var itemsgb []list.Item
var updatedMap = make(map[string]string)
var sprawdzam []string
var howgb = 0
var sprawdzanieczegos int
var mapConvert []string

func removeFirstAlphanumeric(s string) string {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+`)
	return re.ReplaceAllString(s, "")
}

func main() {

	// notatka ze nastawienie na poczatku isgb na false nie dziala ale chyba trzeba cos wlasnie probowac w ta strone
	howgb = 0

x2:
	// dobra czyli jak tutaj nastawie se isgb na false to sie wszystko w ogole psuje i jest tak jak przed zaimplementowaniem tego
	// bo teraz zeby naprawic wiekszosc problemow z tym to powinno jakos na poczatku sie to robic na false zeby za kazdym razem tak tego nie bral
	//isgb = false

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

	// tutaj sprobuje ta mapke zmienic na tablice zeby te kanaly tak nie przeskakiwaly w kolejnosci co chwila

	/*for t := range maps.Keys(channels) {
		fmt.Println(t)
		mapConvert = append(mapConvert, t)
	}*/

	// dobra czyli to nie zadziala dlatego ze po pierwsze w tablicy sa same klucze i nie ma wartosci ich
	// also zeby to faktycznie zadzialalo to mozna po prostu sprobowac przerobic i wartosci i klucze na mapke a potem to jakos laczyc ze soba
	// i ogolnie to zeby to zawsze mialo ta sama kolejnosc to ta mapka powinna byc posortowana a w go sie tak nie da

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

	// ten problem tutaj jest najprawdopodobniej przez to ze w pliku history sie robi jakis whitespace nie wiem czemu
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// tutaj trzeba zrobic if statement ze jak isHistory jest na false to sie robi to a jak nie to inaczej
	if !isHistory {
		ator, _ := soup.Get(link)
		atorvid := soup.HTMLParse(ator)
		// tutaj zeby to segmentation naprawic to trzeba jakis if statement na dole zrobic ze jak nie dziala to wychodzi normalnie
		owtput := atorvid.Find("div", "class", "pure-g").HTML()
		re := regexp.MustCompile(`<a[^>]*>.*?<p dir="auto">.*?</p>.*?</a>`)
		matches := re.FindAllString(owtput, -1)
		for _, match := range matches {

			// usuwanie elementow htmla
			match = strings.Replace(match, "<a href=", "", -1)
			if video < 10 {
				match = strings.Replace(match, `><p dir="auto">`, "			", -1)
			} else {
				match = strings.Replace(match, `><p dir="auto">`, "       ", -1)
			}
			//match = strings.Replace(match, `><p dir="auto">`, "			", -1)
			match = strings.Replace(match, `</p></a>`, "", -1)
			// usuwanie cudzyslowiow
			match = strings.Replace(match, `"`, "", -1)
			// na dole mi zamienia link z samych znaczkow na wlasciwy link do youtuba
			match = strings.Replace(match, `/`, fmt.Sprintf("https://www.youtube.com/"), -1)
			// zaznacza wszystko co zaczyna sie na https
			cutter := regexp.MustCompile(`https?://[^\s]+`)
			// znajduje to co sie zaczyna, wywala link do youtuba
			link := cutter.FindString(match)
			match = strings.Replace(match, `https://www.youtube.com/`, fmt.Sprintf(""), -1)
			match = strings.Replace(match, "watch?v=", "", 1)
			// koniec wywalania linku do youtuba
			//cutterlist := regexp.MustCompile(`^?list=[^ ]+ `)
			// na dole ostatnia zmiana w projekcie na dzisiaj
			//match = cutterlist.ReplaceAllString(match, "test")

			// to mi zamienia match na playliste

			testlista := regexp.MustCompile(`\?list=.*?\s`)
			match = testlista.ReplaceAllString(match, "[Playlist]")

			if strings.Contains(match, "[Playlist]") {

				match = strings.ReplaceAll(match, "[Playlist]", "")
				// takie cos dziala normalnie
				match = fmt.Sprintf("%s [Playlist]", match)
				// to ze wyswietlanie historii nagle nie dziala to nie jest wina tego jak cos i w poprzednim pushu z gh to juz tez nie dzialalo tylko tego nie zauwazylem

			}

			/*if strings.Contains(match, "[Playlist]") {
				match = strings.Replace(match, "[Playlist]", "", -1)
				match = fmt.Sprintf("%s [Playlist]", match)
				fmt.Println(match)
				// aha czyli jak to zrobie to playlisty sie w ogole nie pokazuja
			} */

			//match = strings.ReplaceAll(match, " ", "")
			// trimspace nic nie daje tylko to replaceall to sprawia
			//match = strings.TrimSpace(match)
			// jak to zrobie to tez odpala pierwsza playliste

			/*testspacje := regexp.MustCompile(`\s+`)
			match = testspacje.ReplaceAllString(match, "") */
			// jak zrobie to na gorze to mi z jakiegos powodu odpala pierwsza playliste

			//dobra to dziala mozna jeszcze jakos wykombinowac zeby normalnie dawalo [playlist] na poczatku zamiast tych tabow jak przy linku
			// ! teoretycznie w ogole zeby naprawic ze sie link od playlisty nie bedzie wyswietlal to mozna zamienic calego stringa na tablice, wywalic element z ?list= na poczatku i potem znowu na stringa
			match = strings.TrimPrefix(match, "?list=")
			//po tym wywaleniu linku regexem w playlistach pierwsze slowo ich sie nie laduje chyba a swoja
			//match = strings.Replace(match, "?list=", "[Playlist]", 1)
			// da sie normalnie usunac ten link z playlisty bo wczesniej juz to robilem tylko musze sobie przypomniec jak
			// o mozna wywalic wszystko co jest po ?list az do spacji

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
		//fmt.Println(itemshist)
		// !!! dobra czyli generalnie to jest problem z tym ze itemshist sie zapetla w ktoryms momencie chyba
		videos = titleLinkMap

		//fmt.Println(len(itemshist))
		//listHeight = len(videos)
		// a moze trzeba czesc tego kodu w ogole przeniesc do maina zamiast trzymac we frameworku
	}

	// pierwszy ifstatement z isgb jest tutaj
	// also moze to jest problem we frameworku trzeba porownac tez z nowa a nie bo w nowej sie tez robi out of range
	if isgb {
		lendeleting := len(mecze)
		howgb += len(mecze)

		// to mi sie moze przydac bo przy opcji wyszukiwania wtedy moge w jakis sposob te wszystkie zliczone filmiki usunac przy wyszukiwaniu
		// also notatka ze to isgb to jest globalna wiec na dole cos z searchem dopiero
		// tutaj pracowalem
		// dobra czyli widze ze teraz generalnie wracanie sie dziala tylko jak juz jest to isgb to play next video i go back sie nie wyswietlaja z jakiegos powodu
		// tutaj tez teraz jest jakis problem z index out of range po jakims czasie zmieniania kanalow z isgb z jakiegos powodu
		// to na gorze jest spowodowane dodaniem isgb w string contains

		// dobra czyli mozliwe ze tutaj jest problem z tym mechanizmem ze usuwa poprzednie filmiki z listy
		// ja tu chyba dawalem jakis ifstatement zeby to naprawic z tego co mi sie przypomina

		// !!dobra widze gdzie jest problem z tym out of range niektore kanaly maja mniej niz 60 filmikow i od tego to sie robi
		// widze ze teraz chyba dziala
		sprawdzam = mecze[len(mecze)-lendeleting:]
		// ale jak cos to tutaj jest caly problem w linijce na gorze
		// tutaj jak cos jest rozwiazanie do tego w nowszej wersji tylko trzeba teraz to wszystko poprzenosic
		// dobra chyba jednak nie ma bo jak teraz odpalam nowa wersje to tez sie robi out of range po jakims czasie
		// w sumie to moze byc problem z tymi pointerami wiec notatka na jutro zeby sprawdzic czy dalej sie bedzie robic jak je wywale nie to nie przez to
		fmt.Println(len(sprawdzam))
		itemsgb = itemsgb[:0]
		// btw na gore kodu nawet nie ma co patrzec w tym przypadku bo tam nie ma nic z isgb poza zadeklarowaniem howgb na zero i sprawdzam jako listy
		// to na gorze jest do wywalania tych filmikow z poprzednich list
		// ciekawe bo tablica sprawdzam ma zawsze 30 elementow jak ja tu wypisuje
		// w sumie jak juz w ogole nie bedzie wyjscia to mozna sprobowac zduplikowac caly ten kod i go przepisac dla isgb
		//fmt.Println(sprawdzam)

		// zmienna sprawdzam jest git ale potem jak ja przypisuje do mecze to juz nie dziala
		// mecze = sprawdzam
	}

	for i, str := range mecze {
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
		mecze[i] = strings.ReplaceAll(str, "&#39;", "")
		mecze[i] = strings.ReplaceAll(str, "&#34;", "")
	}

	for i := range mecze {
		// z tym bledem na dole to cos tutaj w sumie moze byc albo w sumie nie bo to goto drugim w mainie i tak mnie cofa do poczatku wiec moze trzeba jakos zrobic zeby przywracalo wartosci framework.go do tych co byly na poczatku
		if isgb {

			//itemstwo = append(itemstwo, item(sprawdzam[i]))
			itemsgb = append(itemsgb, item(sprawdzam[i]))

		}
		itemstwo = append(itemstwo, item(mecze[i]))

		// tylko mecze to jest w ogole ta tablica wiec w teorii ten kod jest w ogole niepotrzebny
		// itemki jest uzywane 3 razy i potem nic sie z tym nie dzieje
		// jednak sie dzieje jak cos bo to jest zrobione do play next video

		itemki = append(itemki, mecze[i])

		//}

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

		m = &modeltwo{list: l}

		if _, err := tea.NewProgram(m.(tea.Model)).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	} else {

		// jak tutaj pozmienialem tego strukture to sie w ogole cos ekstremalnie popsulo bo w ogole sie listy nie wyswietlaly
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
			// isgb = false tutaj to niby za pierwszym razem dziala ale potem sie robi to samo co wczesniej
		} else {
			// to jest w ogole ciekawe bo na gorze jest isgb wiec on tutaj chyba zaklada ze to nie jest isgb dlatego moze to sie nie wykonuje
			if strings.Contains(link, "search") {

				// jak po tym dam play next video to chyba sie wyswietla z jakiejs innej listy filmik
				//fmt.Println("test test test test")
				//fmt.Println(howgb)
				// howgb += len(mecze) to nie zadziala bo mecze to sa te mecze ktore sie robia od szukania
				//fmt.Println(howgb)
				// wychodzi na to ze howgb tutaj gdzies sie zeruje

				// ale dobra z tego co jest na gorze wynika ze tutaj to normalnie wylapuje jak jest search

				// tutaj jest problem ze jak sie cofne minimum raz to sie nie pokazuja wszystkie filmiki z searcha
				// notatka ze tutaj sie robi caly ten problem z historia ostatni i to sie prawdopodobnie robi dlatego ze dalem wczesniej isgb na false zeby naprawic buga z tym ze po historii search nie dziala

				// poza dodawaniem przy ifstatemencie tego howgb  to to jest jedyny moment w ktorym to jest faktycznie uzywane
				// a to na dole jest do wywalania tych filmikow wszystkich co sie pododawaly wczesniej z searcha

				// widze ze to tutaj jest zrobione normalnie jak powinno byc tylko z jakiegos powodu nie dziala albo tego nie lapie
				sprawdzanieczegos = len(mecze)
				// czyli len od mecze jest zrobione normalnie ze to by dzialalo tylko trzeba ogarnac teraz usuwanie tego

				// dobra czyli z tego co przetestowalem to wynika ze to na dole sie w ogole nie wykonuje
				// itemkidwa := itemstwo
				// fmt.Println(itemstwo)

				itemstwo = itemstwo[howgb:]
				finalremove := len(itemstwo) - len(mecze)
				itemstwo = itemstwo[finalremove:]
				// z dwukropkiem na koncu usunal pierwsze 17 a na poczatku wszystko poza pierwszymi 17
				//fmt.Println(itemstwo)
				// zrobienie nowej zmiennej tez nie dziala wiec tutaj problem lezy w ucinaniu tego a nie w samej zmiennej / liscie
				//fmt.Println(itemstesty)
				//czyli po wyprintowaniu tego 2 razy wychodzi na to ze to uciecie itemstwo nic tutaj nie daje w ogole

				// notatka ze wywalanie rzeczy z listy inna metoda niz dwukropkiem nie dziala bo to nie tablica i ma inny typ
				// trzeba sprobowac zrobic dla searcha nowa liste chyba
				if howgb > 0 {
					//itemstwo = itemstwo[sprawdzanieczegos:]
					fmt.Println(itemstwo)
					//itemstwo = itemstwo[:len(mecze)]
				}

				itemki = itemki[howgb:]
				itemki = itemki[finalremove:]
				// fmt.Println(itemki)
				// czyli przy samym play next video to nie jest problem z lista itemki bo ona to sa normalnie dobre filmiki wszystkie jak powinno byc

				// teraz widze ze play next video jakos dziala ale nie wiem dokladnie o co z tym chodzi bo to w ogole wyswietla chyba jakis losowy filmik z tej listy wiec musze ogarnac mniej wiecej jak to dziala
			}

			//itemstwo = itemstwo[:sprawdzanieczegos]
			//fmt.Println(itemstwo)
			fmt.Println(len(itemstwo))
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
		}
	}

	// sprobuje jutro wyklikac ze wiecej filmikow bo moze one dlatego sie zapetlaja
	/*
		dobra czyli ogolnie kolejnosc jest zalatwiona tylko teraz trzeba ogarnac to zapetlanie i cala funkcja historii bedzie zrobiona
	*/

	if isgb {
		// czyli ze jesli isgb jest nastawione na tak co sie robi zawsze jak kilkne q zeby zmienic kanal to ten kod na dole sie w ogole nie bedzie executowal przez goto statement
		// moze to isgb trzeba przeniesc jakos do srodka funkcji zeby to nie byla globalna zmienna to wtedy cos sie uda
		// przeniesienie tego na dol nic nie daje jak cos bo wtedy ta lista z go back play next itp sie wyswietla nawet jak nic nie wybralem

		// jak cos to to nie jest problem z wartoscia linku tylko z lista prawdopodobnie
		// also w ogole teraz jak to naprawilem to play next video ekran sie nie odpala tylko wraca do wyboru filmiku a jak za pierwszym razem odpalam bez cofania to normalnie jest
		//link = ""
		// isgb = false tutaj tez sie robi to samo co wczesniej
		goto x2

		// !! o czyli jak sie wraca to sie wyswietlaja normalnie filmiki z innych kanalow tylko na poczatku tablicy a potem z tego co byl potem wybrany
		// executowanie binarki na nowo nie dziala jak cos
		// w ogole mozna cos poprobowac z tymi loopami jak mowili zamiast goto
		// czyli tutaj generalnie najlepiej by bylo zeby to restartowalo caly program zamiast robic to goto
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
		goto x
	}
	os.Remove("output.txt")
}
