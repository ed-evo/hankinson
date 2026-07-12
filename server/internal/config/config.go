package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Path to the frontend directory where FE build files are located
	FrontendBasePath string `env:"FE_BASE_PATH" envDefault:"public"`
}

var (
	appConfig *Config
	once      sync.Once
)

func Get() *Config {
	once.Do(func() {
		appConfig = &Config{}

		if err := cleanenv.ReadConfig(".env", appConfig); err != nil {
			log.Fatalf("FATAL: Cannot read config: %v", err)
		}
	})
	return appConfig
}
