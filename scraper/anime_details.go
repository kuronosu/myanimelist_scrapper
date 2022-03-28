package scraper

import (
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

type AnimeDetails struct {
	Name  string `json:"name"`
	Image string `json:"image"`

	Type        string   `json:"type"`
	Episodes    int      `json:"episodes"`
	Status      string   `json:"status"`
	Aired       string   `json:"aired"`
	Premiered   string   `json:"premiered"`
	Broadcast   string   `json:"broadcast"`
	Producers   []string `json:"producers"`
	Licensors   []string `json:"licensors"`
	Studios     []string `json:"studios"`
	Source      string   `json:"source"`
	Genres      []string `json:"genres"`
	Theme       string   `json:"theme"`
	Demographic string   `json:"demographic"`
	Duration    string   `json:"duration"`
	Rating      string   `json:"rating"`

	Score      float32 `json:"score"`
	ScoredBy   int     `json:"scored_by"`
	Ranked     int     `json:"ranked"`
	Popularity int     `json:"popularity"`
	Members    int     `json:"members"`
	Favorites  int     `json:"favorites"`

	StreamingPlatforms []string `json:"streaming_platforms"`
}

func buildAnimeDetails(name, image string, lse LeftsideElements, streamingPlatforms []string) AnimeDetails {
	score, scoredBy := lse.Score()
	return AnimeDetails{
		Name:               name,
		Image:              image,
		Type:               lse.Type(),
		Episodes:           lse.Episodes(),
		Status:             lse.Status(),
		Aired:              lse.Aired(),
		Premiered:          lse.Premiered(),
		Broadcast:          lse.Broadcast(),
		Producers:          lse.Producers(),
		Licensors:          lse.Licensors(),
		Studios:            lse.Studios(),
		Source:             lse.Source(),
		Genres:             lse.Genres(),
		Theme:              lse.Theme(),
		Demographic:        lse.Demographic(),
		Duration:           lse.Duration(),
		Rating:             lse.Rating(),
		Score:              score,
		ScoredBy:           scoredBy,
		Ranked:             lse.Ranked(),
		Popularity:         lse.Popularity(),
		Members:            lse.Members(),
		Favorites:          lse.Favorites(),
		StreamingPlatforms: streamingPlatforms,
	}
}

type safeAnimeContainer struct {
	mu     sync.Mutex
	animes map[string]AnimeDetails
}

func (container *safeAnimeContainer) AddAnime(url string, anime AnimeDetails) {
	container.mu.Lock()
	defer container.mu.Unlock()
	container.animes[url] = anime
}

func ScrapeAnimes(urls []string) map[string]AnimeDetails {
	c := colly.NewCollector(colly.Async(true))
	container := safeAnimeContainer{
		animes: make(map[string]AnimeDetails),
	}

	c.OnError(func(r *colly.Response, err error) {
		log.Println("\033[31mSomething went wrong\033[0m: ", r.Request.URL, err)
		log.Fatal()
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("\033[32mVisited\033[0m", r.Request.URL)
	})

	c.OnHTML("div#contentWrapper", func(e *colly.HTMLElement) {
		// fmt.Println("Name: ", getName(e))
		// fmt.Println("Image: ", getImage(e))
		// fmt.Println("StreamingPlatforms: ", getStreamingPlatforms(e))
		anime := buildAnimeDetails(getName(e), getImage(e), getLeftsideElements(e), getStreamingPlatforms(e))
		container.AddAnime(e.Request.URL.String(), anime)
	})

	for _, url := range urls {
		c.Visit(url)
	}
	c.Wait()
	return container.animes
}

func getName(contentWrapper *colly.HTMLElement) string {
	return contentWrapper.ChildText(".title-name.h1_bold_none")
}

func getImage(contentWrapper *colly.HTMLElement) string {
	// return leftside.ChildAttr("img[itemprop]", "src") //data-src
	return contentWrapper.ChildAttr("img[itemprop]", "data-src") //data-src
}

func getStreamingPlatforms(contentWrapper *colly.HTMLElement) []string {
	streamingPlatforms := make([]string, 0)
	contentWrapper.ForEach("div.broadcast", func(i int, el *colly.HTMLElement) {
		streamingPlatforms = append(streamingPlatforms, strings.TrimSpace(el.Text))
	})
	return streamingPlatforms
}
