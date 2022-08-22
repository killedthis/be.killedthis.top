package builder

import (
	"fmt"
	"github.com/ryanbradynd05/go-tmdb"
	"io"
	"killedthis/shared"
	"log"
	"net/http"
	"os"
)

const TmdbBaseUrl = "https://image.tmdb.org/t/p/w300_and_h450_bestv2"

type media []struct {
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title"`
	ID          string `json:"ID"`
	PosterPath  string `json:"poster_path"`
}

type TmdbPosterDownloader struct {
	config tmdb.Config
}

func NewPosterDownloader(tmdbApiKey string) *TmdbPosterDownloader {
	if tmdbApiKey == "" {
		log.Panic("TMDB_API not defined")
		return nil
	}

	return &TmdbPosterDownloader{
		config: tmdb.Config{
			APIKey:   string(tmdbApiKey),
			Proxies:  nil,
			UseProxy: false,
		},
	}
}

func (m TmdbPosterDownloader) LookupPosters(baseFolder string, shows []shared.KilledShow) {
	// Check if base folder exists
	if _, err := os.Stat(baseFolder); err != nil {
		err = os.MkdirAll(baseFolder, os.ModePerm)
		if err != nil {
			log.Panic("failed to create posters base folder: ", err)
			return
		}
	}

	tmdbAPI := tmdb.Init(m.config)

	for _, show := range shows {
		// skip if no tmdb id defined
		if show.TmdbId == nil {
			continue
		}

		// Check if poster already exists
		posterFilename := fmt.Sprintf("%s/%d.jpg", baseFolder, *show.TmdbId)
		log.Printf("Checking poster '%s'...\n", posterFilename)

		if _, err := os.Stat(posterFilename); err == nil {
			log.Println("\t...exists")
		}

		// if not, look up the path using the tmdb api
		log.Printf("Poster for show '%s' not found. Looking up URL on TMDB...", show.Title)
		tvInfo, err := tmdbAPI.GetTvInfo(*show.TmdbId, nil)
		if err != nil {
			log.Println("Failed to query TMDB API, skippig: ", err)
			continue
		}

		// download poster
		m.downloadPoster(posterFilename, TmdbBaseUrl+tvInfo.PosterPath)
	}
}

func (m TmdbPosterDownloader) downloadPoster(posterFilename string, url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Println("failed to retrieve poster: ", err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		file, err := os.Create(posterFilename)
		if err != nil {
			log.Println("failed to create poster file: ", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Println("failed to sync file: ", err)
			return
		}
	}
}
