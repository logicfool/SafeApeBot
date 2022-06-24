package bot

import (
	"SafeApeBot/db"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) SendText(chatid interface{}, text string) {
	var chatId int64
	switch chatid := chatid.(type) {
	case int64:
		chatId = chatid
	case string:
		chatId, _ = strconv.ParseInt(chatid, 10, 64)
	}
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Bot.Send(msg)
}

func (bot *Bot) EditMessageText(old_message tgbotapi.Message, text string, Keyboard interface{}) {
	msg := tgbotapi.NewEditMessageText(old_message.Chat.ID, old_message.MessageID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if Keyboard != nil {
		switch Keyboard := Keyboard.(type) {
		case tgbotapi.InlineKeyboardMarkup:
			key := Keyboard
			msg.ReplyMarkup = &key
		case tgbotapi.ReplyKeyboardMarkup:
			key := Keyboard
			_ = key
		}
	}
	bot.Bot.Send(msg)
}

func (bot *Bot) DeleteMessage(message tgbotapi.Message) {
	msg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	bot.Bot.Send(msg)
}

func (bot *Bot) AddBotToRoot(update tgbotapi.Update) {
	curr_bot := bot.DB.GetBotByToken(bot.Bot.Token)
	if curr_bot.BotToken == "" {
		curr_bot.BotToken = bot.Bot.Token
		curr_bot.RootBot = true
		curr_bot.Enabled = true
		curr_bot.SudoUser = update.Message.From.ID
		if curr_bot.Chats == nil {
			curr_bot.Chats = make(map[int64]db.DBChatStruct)
		}
		bot.DB.AddBot(curr_bot)
	} else {
		if curr_bot.Chats == nil {
			curr_bot.Chats = make(map[int64]db.DBChatStruct)
		}
		curr_bot.RootBot = true
		bot.DB.UpdateBot(curr_bot)
	}
}

func (bot *Bot) AddBot(update tgbotapi.Update) {
	curr_bot := bot.DB.GetBotByToken(bot.Bot.Token)
	if curr_bot.BotToken == "" {

		curr_bot.BotToken = bot.Bot.Token
		curr_bot.RootBot = false
		curr_bot.Enabled = true
		curr_bot.SudoUser = update.Message.From.ID
		if curr_bot.Chats == nil {
			curr_bot.Chats = make(map[int64]db.DBChatStruct)
		}
		err := bot.DB.AddBot(curr_bot)
		if err != nil {
			bot.SendText(update.Message.Chat.ID, "Error: "+err.Error())
		}
		bot.SendText(update.Message.Chat.ID, "Bot added")
	} else {
		if curr_bot.Chats == nil {
			curr_bot.Chats = make(map[int64]db.DBChatStruct)
		}
		bot.DB.UpdateBot(curr_bot)
	}
}
