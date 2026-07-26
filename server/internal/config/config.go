package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GeminiApiKey string `env:"GEMINI_API_KEY"`
	GeminiModel  string `env:"GEMINI_MODEL"`
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
