package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configStore struct {
	open_ai_api_key string
}

var config configStore
var loaded bool = false

func Load() configStore {

	if loaded {
		return config
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := configStore{
		open_ai_api_key: os.Getenv("OPENAI_API_KEY"),
	}
	loaded = true
	return config

}
