package config

type Config struct {
	Api
	Cors
	Database
}

func New() *Config {
	return &Config{
		Api:      API(),
		Cors:     NewCors(),
		Database: DataStore(),
	}
}
