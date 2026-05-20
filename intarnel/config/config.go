package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Panicln("Warning: .env file not found envirement variables!")
	}
	var config *Config = &Config{
		DatabaseURL: os.Getenv("DATABASEURL"),
		Port:        os.Getenv("PORT"),
	}
	return config, nil
}
