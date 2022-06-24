package bothook

import (
	"SafeApeBot/bot"
	"SafeApeBot/structs"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Bot struct {
		Bot        *tgbotapi.BotAPI
		Token      string
		Update     tgbotapi.Update
		W          http.ResponseWriter
		BotSession bot.Bot
		Config     *structs.SConfig
	}
)

func InitFromTokenAndUpdate(token string, update tgbotapi.Update, writer interface{}) (Bot, error) {
	bot := Bot{Token: token, Update: update}
	err := bot.InitBot()
	if err != nil {
		return bot, err
	}
	if writer != nil {
		bot.W = writer.(http.ResponseWriter)
	}
	return bot, nil
}

func InitFromConfigAndUpdate(config *structs.SConfig, update tgbotapi.Update, writer interface{}) (Bot, error) {
	bot := Bot{Config: config, Update: update}
	err := bot.InitBot()
	if err != nil {
		return bot, err
	}
	if writer != nil {
		bot.W = writer.(http.ResponseWriter)
	}
	return bot, nil
}
func (bots *Bot) InitBot() error {
	var err error
	// bott, err := tgbotapi.NewBotAPI(bots.Token)
	if bots.Config != nil {
		bots.Token = bots.Config.BotToken
		bots.BotSession = bot.InitFromConfig(bots.Config)
		bots.Bot = bots.BotSession.Bot
		return nil
	}
	bots.BotSession, err = bot.InitByToken(bots.Token)
	if err != nil {
		return err
	}
	bots.Bot = bots.BotSession.Bot
	// bots.BotSession.Config = bots.Config
	return nil
}

func (bot *Bot) HandleUpdate() {
	if bot.Bot == nil {
		bot.InitBot()
	}
	// spew.Dump(bot.Update.Message)
	// return
	if bot.Update.Message != nil {
		text := bot.Update.Message.Text
		if len(text) >= 1 {
			if bot.Update.Message.Text[0] == '/' {
				bot.handleCommand()
				return
			} else {
				bot.BotSession.HandleMessage(bot.Update)
			}
		}
	} else {
		//might be an edited message!
		// spew.Dump(bot.Update)
		bot.BotSession.HandleUpdate(bot.Update)
		return
	}

}

func (bot *Bot) handleCommand() {
	// if len > 0
	text := bot.Update.Message.Text
	_ = text

	// command is the first word split by space
	command := strings.Split(text[1:], " ")[0]
	_ = command
	switch command {
	case "start":
		msg := bot.BotSession.GetStartMessage(bot.Update.Message.Chat.ID, false)
		bot.BotSession.Bot.Send(msg)
	case "help":
		msg := bot.BotSession.GetHelpMessage(bot.Update.Message.Chat.ID, false)
		bot.BotSession.Bot.Send(msg)
	default:
		bot.BotSession.HandleCommand(bot.Update)
	}
	// fmt.Println("Command: ", command)
	// fmt.Fprint(bot.W, "Command: ", command)
}

func (bot *Bot) handleText() {
	text := bot.Update.Message.Text
	_ = text

	bot.BotSession.HandleMessage(bot.Update)
}
