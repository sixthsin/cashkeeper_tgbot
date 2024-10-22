package telegram

import (
	"context"
	"gotgbot/pkg/storage/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCommand(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	waitingUsers map[int64]chan string,
	s *db.Storage,
) {
	_, isWaiting := waitingUsers[update.Message.From.ID]
	switch {
	case update.Message.Text == StartCmd: // обработка команды /start
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgStart)
		bot.Send(msg)

	case update.Message.Text == HelpCmd: // обработка команды /help
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgHelp)
		bot.Send(msg)

	case update.Message.Text == AddCategoryCmd: // обработка команды /add_category
		dataChan := make(chan string)
		waitingUsers[update.Message.From.ID] = dataChan
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgEnterCatrgoryData)
		bot.Send(msg)

		data := <-dataChan
		go addCategory(ctx, bot, update, s, data)

		delete(waitingUsers, update.Message.From.ID)

	case update.Message.Text == DeleteCategoryCmd: // обработка команды /delete_category
		dataChan := make(chan string)
		waitingUsers[update.Message.From.ID] = dataChan
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgCategoryDelete)
		bot.Send(msg)

		data := <-dataChan
		go deleteCategory(ctx, bot, update, s, data)

		delete(waitingUsers, update.Message.From.ID)

	case update.Message.Text == GetCategoriesListCmd: // обработка команды /my_categories
		go getCategoriesList(ctx, bot, update, s)

	case update.Message.Text == AddExpensesCmd: // обработка команды /add_expenses
		dataChan := make(chan string)
		waitingUsers[update.Message.From.ID] = dataChan

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgEnterExpensesData)
		bot.Send(msg)

		data := <-dataChan
		go addExpenses(ctx, bot, update, s, data)

	case isWaiting && update.Message.Text == BackCmd: // обработка команды /back
		delete(waitingUsers, update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgDeny)
		bot.Send(msg)

	case isWaiting:
		waitingUsers[update.Message.From.ID] <- update.Message.Text

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsgUnknownCmd)
		bot.Send(msg)
	}
}
