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
		renderer.RenderHtml()
	}
}
