package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig interface {
	LoadEnv()
}

type appConfig struct {
}

func (cfg appConfig) LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		os.Exit(1)
	}
	log.Println("Successfully loaded .env file")
}

func New() AppConfig {
	return &appConfig{}
}
