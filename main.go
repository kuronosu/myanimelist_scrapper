package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kuronosu/myanimelist_scrapper/scraper"
)

type DataContainer struct {
	TopAnimes []scraper.TopAnime     `json:"top_animes"`
	Animes    []scraper.AnimeDetails `json:"animes"`
}

func SaveDataContainer(filename string, data DataContainer) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("animes.json", b, 0644)
}

func LoadDataContainer(filename string) (DataContainer, error) {
	file, _ := ioutil.ReadFile(filename)
	data := DataContainer{}
	err := json.Unmarshal([]byte(file), &data)
	return data, err
}

// if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {

// }

func getDataContainerFromNetwork() DataContainer {
	return DataContainer{
		TopAnimes: scraper.ScrapeTopAnimes(1),
		Animes:    []scraper.AnimeDetails{},
	}
}

func main() {
	// animes := scraper.ScrapeTopAnimes(1)
	data := getDataContainerFromNetwork()
	SaveDataContainer("top-animes.json", data)
}
