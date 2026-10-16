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

	DATABASE_URL string

	// Auth
	AUTH_OAUTH_COOKIE_SECRET string
	AUTH_APP_COOKIE_SECRET   string

	// ioprodz account (Zitadel OIDC) and the services portal
	IOPRODZ_AUTH_ISSUER        string
	IOPRODZ_AUTH_CLIENT_ID     string
	IOPRODZ_AUTH_CLIENT_SECRET string
	PORTAL_SERVICES            string

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

		IS_PRODUCTION:              isProduction,
		ENVIRONMENT:                os.Getenv("APP_ENV"),
		PORT:                       os.Getenv("PORT"),
		APP_DOMAIN:                 os.Getenv("APP_DOMAIN"),
		DATABASE_URL:               os.Getenv("DATABASE_URL"),
		AUTH_OAUTH_COOKIE_SECRET:   os.Getenv("AUTH_OAUTH_COOKIE_SECRET"),
		AUTH_APP_COOKIE_SECRET:     os.Getenv("AUTH_APP_COOKIE_SECRET"),
		IOPRODZ_AUTH_ISSUER:        getEnvOr("IOPRODZ_AUTH_ISSUER", "https://auth.ioprodz.com"),
		IOPRODZ_AUTH_CLIENT_ID:     os.Getenv("IOPRODZ_AUTH_CLIENT_ID"),
		IOPRODZ_AUTH_CLIENT_SECRET: os.Getenv("IOPRODZ_AUTH_CLIENT_SECRET"),
		PORTAL_SERVICES:            os.Getenv("PORTAL_SERVICES"),
		OPEN_AI_API_KEY:            os.Getenv("OPENAI_API_KEY"),
		BASE_URL:                   getBaseUrl(isProduction, os.Getenv("PORT"), os.Getenv("APP_DOMAIN")),
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

func getEnvOr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
