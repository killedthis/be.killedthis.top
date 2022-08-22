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

type ConfigurationRoot struct {
	Database        DatabaseConfig `yaml:"database"`
	OutputDirectory string         `yaml:"outputDirectory"`
	TmdbApi         string         `yaml:"tmdbApi"`
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
