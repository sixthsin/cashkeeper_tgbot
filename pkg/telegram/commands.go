package telegram

import (
	"context"
	"fmt"
	"gotgbot/pkg/storage/db"
	"log"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendError(
	err error,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) {
	log.Printf("Error adding user: %v", err)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgDefault+errMsgDeny)
	bot.Send(msg)
}

func addCategory(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	s *db.Storage,
	data string,
) {
	dataBufer := strings.Split(data, " ")
	total, err := strconv.Atoi(dataBufer[len(dataBufer)-1])
	categoryTotalType := reflect.TypeOf(total).Kind()

	if err != nil {
		log.Println("can't convert data: ", err)
		sendError(err, bot, update)
		return
	}

	if len(dataBufer) < minUserMessageLen || categoryTotalType != reflect.Int {
		err := fmt.Errorf("wrong data format")
		sendError(err, bot, update)
		return
	}

	if total < minTotalValue {
		err := fmt.Errorf("wrong category total")
		sendError(err, bot, update)
		return
	}

	category := &db.Category{
		UserID: update.Message.From.ID,
		Title:  strings.Join(dataBufer[:len(dataBufer)-1], " "),
		Total:  total,
	}

	exists, err := s.IsExists(ctx, category)
	if err != nil {
		log.Println("can't chek is category exists: ", err)
		sendError(err, bot, update)
		return
	}

	if exists {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgAlreadyUsed+errMsgDeny)
		bot.Send(msg)
		return
	}

	err = s.Save(ctx, category)

	if err != nil {
		sendError(err, bot, update)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgSuccess)
	bot.Send(msg)
}

func deleteCategory(ctx context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	s *db.Storage,
	data string,
) {
	category := &db.Category{
		UserID: update.Message.From.ID,
		Title:  data,
	}

	isExists, err := s.IsExists(ctx, category)
	if err != nil {
		log.Println("can't chek is category exists: ", err)
		sendError(err, bot, update)
		return
	}

	if !isExists {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgUnknownTitle+errMsgDeny)
		bot.Send(msg)
		return
	}

	err = s.Delete(ctx, category)

	if err != nil {
		sendError(err, bot, update)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgSuccess)
	bot.Send(msg)
}

func getCategoriesList(ctx context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	s *db.Storage,
) {
	var msgText string
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgCategoriesList)
	bot.Send(msg)
	rowsList, err := s.Get(ctx, update.Message.From.ID)
	if err != nil {
		sendError(err, bot, update)
		return
	}
	for _, row := range rowsList {
		remained := row.Total - row.TotalSpent
		msgText += "\"" + row.Title + "\":\n" + "Лимит - " + strconv.Itoa(row.Total) + "\nПотрачено - " + strconv.Itoa(row.TotalSpent) + "\nОсталось - " + strconv.Itoa(remained) + "\n\n"
	}
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}

func addExpenses(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	s *db.Storage,
	data string,
) {
	dataBufer := strings.Split(data, " ")
	totalSpent, err := strconv.Atoi(dataBufer[len(dataBufer)-1])
	categoryTotalSpentType := reflect.TypeOf(totalSpent).Kind()
	if err != nil {
		log.Println("can't convert data: ", err)
		sendError(err, bot, update)
		return
	}

	if len(dataBufer) < minUserMessageLen || categoryTotalSpentType != reflect.Int {
		err := fmt.Errorf("wrong data format")
		sendError(err, bot, update)
		return
	}

	category := &db.Category{
		UserID:     update.Message.From.ID,
		Title:      strings.Join(dataBufer[:len(dataBufer)-1], " "),
		TotalSpent: totalSpent,
	}

	isExists, err := s.IsExists(ctx, category)
	if err != nil {
		log.Println("can't chek is category exists: ", err)
		sendError(err, bot, update)
		return
	}

	if !isExists {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgUnknownTitle+errMsgDeny)
		bot.Send(msg)
		return
	}

	err = s.Update(ctx, category)

	if err != nil {
		sendError(err, bot, update)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgSuccess)
	bot.Send(msg)
}
