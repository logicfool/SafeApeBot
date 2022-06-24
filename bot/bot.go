package bot

import (
	"SafeApeBot/db"
	"SafeApeBot/structs"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Bot struct {
		Bot            *tgbotapi.BotAPI
		Me             tgbotapi.User
		Config         *structs.SConfig
		DB             *db.DB
		CommandHandler *CommandHandler
	}
)

func InitFromConfig(config *structs.SConfig) Bot {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {

	}
	// bot.Debug = true
	bott := Bot{Bot: bot, Me: tgbotapi.User{}, CommandHandler: &CommandHandler{}, Config: config}
	return bott
}

func InitByToken(token string) (Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return Bot{}, err
	}
	// myself, _ := bot.GetMe()
	// bot.Debug = true
	bott := Bot{Bot: bot, Me: tgbotapi.User{}, CommandHandler: &CommandHandler{}}
	return bott, nil
}

func (bot *Bot) SetConfig(config *structs.SConfig) {
	bot.Config = config
	bot.AddCommands()
}

func (bot *Bot) GetMe() {
	bot.Me, _ = bot.Bot.GetMe()
}
func (bot *Bot) Start() {
	bot.Me, _ = bot.Bot.GetMe()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.Bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		// spew.Dump(update)
		if update.Message != nil {
			bot.HandleMessage(update)
		} else if update.CallbackQuery != nil {
			bot.handleCallbackQuery(update)
		}
	}
}

func (bot *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.NewChatMembers != nil {
			bot.handleNewMemberJoin(update)
			return
		}
		bot.HandleMessage(update)
	} else if update.CallbackQuery != nil {
		bot.handleCallbackQuery(update)
	} else if update.MyChatMember != nil {
		return
	}
}

func (bot *Bot) HandleMessage(update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}
	isservicemessage := bot.IsServiceMessage(update)
	if isservicemessage {
		bot.DeleteMessage(*message)
		return
	}
	if message.IsCommand() {
		bot.HandleCommand(update)
		return
	}
	bot.SearchAdnSendFilter(update)
	// bot.SendText(message.Chat.ID, message.Text)
	// message.Chat.ID
}

func (bot *Bot) handleCallbackQuery(update tgbotapi.Update) {
	callback := update.CallbackQuery
	data := callback.Data
	split := strings.Split(data, ":")[0]
	switch split {
	case "action":
		bot.ActionHandler(update)
	case "help":
		bot.Helphandler(update)
	}
}
