package scraper

import (
	"github.com/gocolly/colly/v2"
)

type AnimeDetails struct {
	Name string `json:"name"`
}

func ScrapeAnime(c colly.Collector, page uint) AnimeDetails {
	return AnimeDetails{}
}
