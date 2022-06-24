package bot

import (
	// bothook "SafeApeBot/server/bothook"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) ActionHandler(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	callbackmessage := update.CallbackQuery.Message
	_, _ = data, callbackmessage
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	_ = msg
	split := strings.Split(data, ":")[1]
	switch split {
	case "help":
		db_data := bot.DB.GetBotByToken(bot.Bot.Token)
		var helpmessage string
		if db_data.RootBot {
			helpmessage = bot.GetHelpMessage(update.CallbackQuery.Message.Chat.ID, true).Text
		} else {
			helpmessage = bot.GetHelpMessage(update.CallbackQuery.Message.Chat.ID, false).Text
		}
		bot.EditMessageText(*callbackmessage, helpmessage, nil)
	}

}

func (bot *Bot) Helphandler(update tgbotapi.Update) {
}
