package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ryanbradynd05/go-tmdb"
)

type Media []struct {
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title"`
	ID          string `json:"ID"`
	PosterPath  string `json:"poster_path"`
}

var PosterPath string = "/home/www/sites/killedthis.top/img/posters/"

func main() {
	var (
		tmdbAPI *tmdb.TMDb
	)

	apikeyload, err := ioutil.ReadFile("/home/www/secret/tmdb")
	if err != nil {
		log.Fatal(err)
	}

	config := tmdb.Config{
		APIKey:   string(apikeyload),
		Proxies:  nil,
		UseProxy: false,
	}

	tmdbAPI = tmdb.Init(config)

	searchPtr := flag.Bool("s", false, "Search Title (string)")
	tvInfoPtr := flag.Bool("t", false, "TV Info (requires id)")

	flag.Parse()

	if *searchPtr == true {
		title := strings.Join(os.Args[2:], " ")
		tvInfo, _ := tmdbAPI.SearchTv(title, nil)

		for md := range tvInfo.Results {
			if title == tvInfo.Results[md].Name {
				// fmt.Println(tvInfo.Results[md].FirstAirDate)
				fmt.Println(tvInfo.Results[md].Name)
				fmt.Println(tvInfo.Results[md].ID)
				// fmt.Println(tvInfo.Results[md].PosterPath)
			}
		}
	}

	if *tvInfoPtr == true {
		baseurl := "https://image.tmdb.org/t/p/w300_and_h450_bestv2"
		id, _ := strconv.Atoi(os.Args[2])
		tvInfo, err := tmdbAPI.GetTvInfo(id, nil)

		if err != nil {
			fmt.Println(os.Args[2:])
			fmt.Println(err)
		}

		// fmt.Println(tvInfo.Overview)
		// fmt.Println(tvInfo.NumberOfEpisodes)
		// fmt.Println(tvInfo.NumberOfSeasons)
		// fmt.Println(tvInfo.FirstAirDate)
		// fmt.Println(tvInfo.LastAirDate)

		// fmt.Println(baseurl + tvInfo.PosterPath)

		downloadPoster(baseurl+tvInfo.PosterPath, os.Args[2])
	}
}

func downloadPoster(url string, tmdbid string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(tmdbid)
		fmt.Println(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		file, err := os.Create(PosterPath + tmdbid + ".jpg")
		if err != nil {
			fmt.Println(tmdbid)
			fmt.Println(err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(tmdbid)
			fmt.Println(err)
		}
	}
}
