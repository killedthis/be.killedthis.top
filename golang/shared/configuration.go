package shared

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type DatabaseConfig struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

type TmdbConfig struct {
	Apikey string `yaml:"apikey"`
}

type ContentGeneratorConfig struct {
	TmdbEnabled     bool   `yaml:"tmdbEnabled"`
	OutputDirectory string `yaml:"outputDirectory"`
	FormatSiteFile  string `yaml:"formatSiteFile"`
}

type ConfigurationRoot struct {
	Database         DatabaseConfig         `yaml:"database"`
	ContentGenerator ContentGeneratorConfig `yaml:"contentGenerator"`
	Tmdb             TmdbConfig             `yaml:"tmdb"`
}

func LoadConfiguration(filename string) *ConfigurationRoot {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		log.Panic("unable to load configuration: ", err)
		return nil
	}

	var config ConfigurationRoot
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Panic("unable to parse configuration: ", err)
		return nil
	}

	return &config
}
