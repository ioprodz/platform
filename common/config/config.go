package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configStore struct {
	OPEN_AI_API_KEY string
	ENVIRONMENT     string
	PORT            string
}

var config configStore
var loaded bool = false

func Load() configStore {

	if loaded {
		return config
	}

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	config := configStore{
		ENVIRONMENT:     os.Getenv("APP_ENV"),
		PORT:            os.Getenv("PORT"),
		OPEN_AI_API_KEY: os.Getenv("OPENAI_API_KEY"),
	}
	loaded = true
	return config

}
