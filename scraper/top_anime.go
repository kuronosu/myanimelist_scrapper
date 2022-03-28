package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type TopAnimeInformation struct {
	Type     string `json:"type"`
	Episodes int    `json:"episodes"`
	Aired    string `json:"aired"`
	Members  int    `json:"members"`
}

type TopAnime struct {
	Url         string              // ✔
	Rank        uint                // ✔
	Name        string              // ✔
	Information TopAnimeInformation // ✔
	Score       float32             // ✔
	// Img         string           // ✘
}

func splitInformation(information string) (string, string, string) {
	tmp := strings.Split(information, "\n")
	return tmp[0], tmp[1], tmp[2]
}

func parseTypeEps(ty_eps string) (string, int) {
	tmp := strings.Split(ty_eps, "(")
	eps, _ := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(tmp[1], "eps)", "")))
	return strings.TrimSpace(tmp[0]), eps
}

func parseMembers(members string) (int, error) {
	members = strings.Replace(members, " members", "", 1)
	members = strings.ReplaceAll(members, ",", "")
	return strconv.Atoi(strings.TrimSpace(members))
}

func parseTopAnimeInformation(data string) TopAnimeInformation {
	ty_eps, emitted, members := splitInformation(data)
	typ, eps := parseTypeEps(ty_eps)
	membersInt, _ := parseMembers(members)
	emitted = strings.TrimSpace(emitted)
	return TopAnimeInformation{
		Type:     typ,
		Episodes: eps,
		Aired:    emitted,
		Members:  membersInt,
	}
}

func ScrapeTopAnimesByPage(page uint) []TopAnime {
	c := colly.NewCollector()
	topAnimeList := make([]TopAnime, 0)

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	c.OnError(func(r *colly.Response, err error) {
		log.Println("\033[31mSomething went wrong\033[0m: ", r.Request.URL, err)
		log.Fatal()
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("\033[32mVisited\033[0m", r.Request.URL)
	})

	c.OnHTML("table.top-ranking-table>tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr.ranking-list", func(i int, el *colly.HTMLElement) {
			rank, err := strconv.ParseUint(el.ChildText(".top-anime-rank-text"), 10, 32)
			if err != nil {
				rank = 0
			}
			name := el.ChildText(".anime_ranking_h3")
			url := el.ChildAttr("a", "href")
			score, err := strconv.ParseFloat(el.ChildText(".js-top-ranking-score-col"), 32)
			if err != nil {
				score = -1
			}
			information := parseTopAnimeInformation(el.ChildText(".information"))

			topAnimeList = append(topAnimeList, TopAnime{
				Rank:        uint(rank),
				Name:        name,
				Url:         url,
				Score:       float32(score),
				Information: information,
			})
		})
	})
	c.Visit("https://myanimelist.net/topanime.php?limit=" + strconv.FormatUint(uint64(page*50), 10))
	c.Wait()
	return topAnimeList
}

func ScrapeTopAnimes(pageCount uint) []TopAnime {
	animes := make([]TopAnime, 0)
	for i := 0; i < int(pageCount); i++ {
		animes = append(animes, ScrapeTopAnimesByPage(uint(i))...)
	}
	return animes
}
