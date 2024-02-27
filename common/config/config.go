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
	BASE_URL      string
	PORT          string
	IS_PRODUCTION bool

	DB_MONGO_URI string

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

	isProduction := os.Getenv("APP_ENV") == "production"
	config = configStore{

		IS_PRODUCTION:            isProduction,
		ENVIRONMENT:              os.Getenv("APP_ENV"),
		PORT:                     os.Getenv("PORT"),
		APP_DOMAIN:               os.Getenv("APP_DOMAIN"),
		DB_MONGO_URI:             os.Getenv("DB_MONGO_URI"),
		AUTH_OAUTH_COOKIE_SECRET: os.Getenv("AUTH_OAUTH_COOKIE_SECRET"),
		AUTH_APP_COOKIE_SECRET:   os.Getenv("AUTH_APP_COOKIE_SECRET"),
		OPEN_AI_API_KEY:          os.Getenv("OPENAI_API_KEY"),
		BASE_URL:                 getBaseUrl(isProduction, os.Getenv("PORT"), os.Getenv("APP_DOMAIN")),
	}

	loaded = true
	return config

}

func getBaseUrl(isProduction bool, port string, appdomain string) string {
	baseUrlProtocol := "http"
	baseUrlPort := ":" + os.Getenv("PORT")
	if isProduction {
		baseUrlProtocol = "https"
		baseUrlPort = ""
	}
	return baseUrlProtocol + "://" + appdomain + baseUrlPort
}
