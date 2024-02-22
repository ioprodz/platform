package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configStore struct {
	// General
	ENVIRONMENT   string
	APP_DOMAIN    string
	PORT          string
	IS_PRODUCTION bool

	// Auth
	AUTH_OAUTH_COOKIE_SECRET string
	AUTH_APP_COOKIE_SECRET   string

	// Services
	OPEN_AI_API_KEY string
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

	config = configStore{

		ENVIRONMENT:              os.Getenv("APP_ENV"),
		IS_PRODUCTION:            os.Getenv("APP_ENV") == "production",
		PORT:                     os.Getenv("PORT"),
		APP_DOMAIN:               os.Getenv("APP_DOMAIN"),
		AUTH_OAUTH_COOKIE_SECRET: os.Getenv("AUTH_OAUTH_COOKIE_SECRET"),
		AUTH_APP_COOKIE_SECRET:   os.Getenv("AUTH_APP_COOKIE_SECRET"),
		OPEN_AI_API_KEY:          os.Getenv("OPENAI_API_KEY"),
	}

	loaded = true
	return config

}
