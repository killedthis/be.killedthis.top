package main

import (
	"killedthis/be.killedthis.top/builder"
	"log"
	"os"
)

func main() {
	log.Println("Starting KilledThis Builder...")

	outputFolder := os.Getenv("OUTPUT")
	if outputFolder == "" {
		log.Panic("unknown output folder, specify ENV 'OUTPUT', should probably go into a config file later")
		return
	}

	database := builder.OpenDatabase()
	if database == nil {
		log.Panic("failed to open database")
	}

	posterDownloader := builder.NewPosterDownloader()

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

		// look up & download posters for shows
		posterDownloader.LookupPosters(outputFolder+"/img/posters", shows)
	}
}
