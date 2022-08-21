package main

import (
	"killedthis/be.killedthis.top/builder"
	"log"
)

func main() {
	log.Println("Starting KilledThis Builder...")

	database := builder.OpenDatabase()
	if database == nil {
		log.Panic("failed to open database")
	}

	// Get all Service Providers
	serviceProviders := database.GetServiceProviders()
	if serviceProviders == nil {
		log.Panic("nothing found")
	}

	// Render the templates for each service provider
	for _, x := range serviceProviders {
		// get the shows killed by this provider
		shows := database.GetShowsByProvider(x)

		// create a new renderer for this provider and shows
		renderer := builder.NewRenderer(x, serviceProviders, shows)

		// render it
		renderer.RenderHtml()
	}
}
