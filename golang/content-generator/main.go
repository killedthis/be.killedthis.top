package main

import (
	"killedthis/content-generator/builder"
	"killedthis/shared"
	"log"
)

var config *shared.ConfigurationRoot

func init() {
	log.Println("Loading Configuration...")

	config = shared.LoadConfiguration("../config.yaml")
	if config == nil {
		log.Panic("failed to load configuration")
	}
}

func main() {
	log.Println("Starting KilledThis Builder...")

	outputFolder := config.ContentGenerator.OutputDirectory
	if outputFolder == "" {
		log.Panic("unknown output folder, specify ENV 'OUTPUT', should probably go into a config file later")
		return
	}

	database := shared.OpenDatabase(&config.Database)
	if database == nil {
		log.Panic("failed to open database")
	}

	var posterDownloader *builder.TmdbPosterDownloader = nil
	if config.ContentGenerator.TmdbEnabled {
		posterDownloader = builder.NewPosterDownloader(config.Tmdb.Apikey)
	}

	// Get all Service Providers
	log.Println("Retrieving service providers...")
	serviceProviders := database.GetServiceProviders()
	if serviceProviders == nil {
		log.Panic("nothing found")
	}

	// Render the templates for each service provider
	for _, provider := range serviceProviders {
		log.Printf("Rendering Page for Provider '%s'...\n", provider)
		// get the shows killed by this provider
		shows := database.GetShowsByProvider(provider)

		// create a new renderer for this provider and shows
		renderer := builder.NewRenderer(provider, serviceProviders, shows)

		// render it
		renderer.RenderHtml(outputFolder)

		if posterDownloader != nil {
			// look up & download posters for shows
			posterDownloader.LookupPosters(outputFolder+"/img/posters", shows)
		}
	}
}
