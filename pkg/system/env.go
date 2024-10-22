package system

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const envPath = "D:/Dev/gotgbot/.env"

func BotToken() string {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Panic(err)
	}
	return os.Getenv("TOKEN")
}

func StoragePath() string {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Panic(err)
	}
	return os.Getenv("STORAGE_PATH")
}
