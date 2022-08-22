package builder

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
