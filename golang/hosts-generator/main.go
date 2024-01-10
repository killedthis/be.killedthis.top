package main

import (
	"killedthis/shared"
	"log"
)

func main() {
	log.Println("Loading Configuration...")

	config := shared.LoadConfiguration("config.yaml")
	if config == nil {
		log.Panic("failed to load configuration")
	}

	database := shared.OpenDatabase(&config.Database)
	if database == nil {
		log.Panic("failed to open database")
	}

	serviceProviders := database.GetServiceProviders()
	if serviceProviders == nil {
		log.Panic("nothing found")
	}

	// Render the templates for each service provider
	for _, provider := range serviceProviders {
		log.Printf("Generating nginx host for Provider '%s'...\n", provider)

	}

	log.Println("Starting KilledThis Host Generator...")

}
