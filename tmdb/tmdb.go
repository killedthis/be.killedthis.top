package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
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
		tmdbid  string
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

		tmdbid = os.Args[2]

		if _, err := os.Stat(PosterPath + tmdbid + ".jpg"); errors.Is(err, os.ErrNotExist) {
			baseurl := "https://image.tmdb.org/t/p/"
			smallImage := "w150_and_h225_bestv2"
			// bigImage := "w780"
			id, _ := strconv.Atoi(tmdbid)
			tvInfo, err := tmdbAPI.GetTvInfo(id, nil)

			if err != nil {
				fmt.Printf("%s | %v\n", tmdbid, err)
			}

			fmt.Printf("Downloading %s from %s\n", tmdbid, baseurl+smallImage+tvInfo.PosterPath)
			downloadPoster(baseurl+smallImage+tvInfo.PosterPath, tmdbid)
		}

		// fmt.Println(tvInfo.Overview)
		// fmt.Println(tvInfo.NumberOfEpisodes)
		// fmt.Println(tvInfo.NumberOfSeasons)
		// fmt.Println(tvInfo.FirstAirDate)
		// fmt.Println(tvInfo.LastAirDate)

		// fmt.Println(baseurl +smallImage + tvInfo.PosterPath)

		if _, err := os.Stat(PosterPath + tmdbid + ".webp"); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Creating webp of %s\n", tmdbid)
			createwebP(tmdbid)
		}
	}
}

func createwebP(tmdbid string) {
	buffer, err := bimg.Read(PosterPath + tmdbid + ".jpg")
	if err != nil {
		fmt.Printf("%s | %v %v\n", tmdbid, os.Stderr, err)
	}

	newImage, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		fmt.Printf("%s | %v\n", tmdbid, err)
	}

	if bimg.NewImage(newImage).Type() == "webp" {
		fmt.Printf("%s | %s\n", tmdbid, "image was converted into webp")
	}

	bimg.Write(PosterPath+tmdbid+".webp", newImage)
}

func downloadPoster(url string, tmdbid string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s | %v\n", tmdbid, err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		file, err := os.Create(PosterPath + tmdbid + ".jpg")
		if err != nil {
			fmt.Printf("%s | %v\n", tmdbid, err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Printf("%s | %v\n", tmdbid, err)
		}
	}
}
