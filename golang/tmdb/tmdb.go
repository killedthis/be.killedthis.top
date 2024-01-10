// / 2>/dev/null ; gorun "$0" "$@" ; exit $?
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	bimg "github.com/h2non/bimg"
	tmdb "github.com/ryanbradynd05/go-tmdb"
)

// type media []struct {
// 	Overview    string `json:"overview"`
// 	ReleaseDate string `json:"release_date"`
// 	Title       string `json:"title"`
// 	ID          string `json:"ID"`
// 	PosterPath  string `json:"poster_path"`
// }

var posterPath = "/home/www/sites/killedthis.top/img/posters/"

func main() {
	var (
		tmdbAPI        *tmdb.TMDb
		tmdbid         string
		baseurl         = "https://image.tmdb.org/t/p/"
		smallImageBase  = baseurl + "w150_and_h225_bestv2"
		largeImageBase  = baseurl + "w780"
	)

	apikeyload, err := os.ReadFile("/home/www/secret/tmdb")
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
	tvInfoPtr := flag.Int("t", -1, "TV Info (requires id)")
	forcePtr := flag.Bool("f", false, "Force Downloading")

	flag.Parse()



	if *searchPtr == true {
		title := strings.Join(os.Args[2:], " ")
		tvInfo, _ := tmdbAPI.SearchTv(title, nil)
		fmt.Printf("Searching for %v\n", title)
		for md := range tvInfo.Results {
			// fmt.Printf("%s [%d]\n", tvInfo.Results[md].Name, tvInfo.Results[md].ID)
			if strings.EqualFold(title, tvInfo.Results[md].Name) {
				// fmt.Println(tvInfo.Results[md].FirstAirDate)
				fmt.Printf("%s [%d]\n", tvInfo.Results[md].Name, tvInfo.Results[md].ID)
				// fmt.Println(tvInfo.Results[md].Name)
				// fmt.Println(tvInfo.Results[md].ID)
				// fmt.Println(tvInfo.Results[md].PosterPath)
			}
		}
	}

	if *tvInfoPtr > 0 {

		tmdbid = strconv.Itoa(*tvInfoPtr)

		_, err := os.Stat(posterPath + tmdbid + ".jpg")
		if err != nil || *forcePtr == true {
			id, _ := strconv.Atoi(tmdbid)

			tvInfo, err := tmdbAPI.GetTvInfo(id, nil)
			if err != nil {
				fmt.Printf("%s | %v\n", tmdbid, err)
			}

			smallImageURL := smallImageBase + tvInfo.PosterPath
			largeImageURL := largeImageBase + tvInfo.PosterPath

			fmt.Printf("Downloading Poster %s from %s\n", tmdbid, smallImageURL)
			downloadPoster(smallImageURL, tmdbid)

			fmt.Printf("Downloading Large Poster %s from %s\n", tmdbid, largeImageURL)
			downloadPoster(largeImageURL, tmdbid+"w780")

		} else {
			fmt.Printf("%s | %s\n", tmdbid, "Poster already exists")
		}

		// fmt.Println(tvInfo.Overview)
		// fmt.Println(tvInfo.NumberOfEpisodes)
		// fmt.Println(tvInfo.NumberOfSeasons)
		// fmt.Println(tvInfo.FirstAirDate)
		// fmt.Println(tvInfo.LastAirDate)

		// fmt.Println(baseurl +smallImage + tvInfo.PosterPath)

		_, err = os.Stat(posterPath + tmdbid + ".webp")
		if err != nil || *forcePtr == true {
			fmt.Printf("Creating webp of %s\n", tmdbid)
			createwebP(tmdbid)
			createwebP(tmdbid + "w780")
		}
	}
}

func createwebP(tmdbid string) {
	buffer, err := bimg.Read(posterPath + tmdbid + ".jpg")
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

	bimg.Write(posterPath+tmdbid+".webp", newImage)
}

func downloadPoster(url string, tmdbid string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s | %v\n", tmdbid, err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		file, err := os.Create(posterPath + tmdbid + ".jpg")
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
