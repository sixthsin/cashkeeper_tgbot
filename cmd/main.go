package main

import (
	"context"
	"gotgbot/pkg/storage/db"
	"gotgbot/pkg/system"
	"gotgbot/pkg/telegram"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	token := system.BotToken()
	sqliteStoragePath := system.StoragePath()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	s, err := db.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	waitingUsers := make(map[int64]chan string)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			ctx := context.Background()
			go telegram.HandleCommand(ctx, bot, update, waitingUsers, s)
		}
	}
}
