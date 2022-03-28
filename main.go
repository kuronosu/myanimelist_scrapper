package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/kuronosu/myanimelist_scrapper/scraper"
)

type DataContainer struct {
	TopAnimes []scraper.TopAnime              `json:"top_animes"`
	Animes    map[string]scraper.AnimeDetails `json:"animes"`
}

func SaveDataContainer(filename string, data DataContainer) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func LoadDataContainer(filename string) (DataContainer, error) {
	file, _ := ioutil.ReadFile(filename)
	data := DataContainer{}
	err := json.Unmarshal([]byte(file), &data)
	return data, err
}

func getDataContainerFromNetwork(page uint) DataContainer {
	topAnimes := scraper.ScrapeTopAnimesByPage(page)
	urls := make([]string, len(topAnimes))
	for i, v := range topAnimes {
		urls[i] = v.Url
	}
	return DataContainer{
		TopAnimes: topAnimes,
		Animes:    scraper.ScrapeAnimes(urls),
	}
}

func ScrapeAndSavePage(page uint) {
	data := getDataContainerFromNetwork(uint(page))
	SaveDataContainer(fmt.Sprintf("data/animes_%d.json", page), data)
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Page to scrape required [0 - n]")
	}
	page, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	data := getDataContainerFromNetwork(uint(page))
	SaveDataContainer(fmt.Sprintf("animes_%d.json", page), data)
	// for i := 0; i < 200; i++ {
	// 	ScrapeAndSavePage(uint(i))
	// 	fmt.Println()
	// }
}
