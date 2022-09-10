package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	CommandHandler struct {
		commands map[string]func(tgbotapi.Update)
	}
)

func (c *CommandHandler) RegisterCommand(cmdstring string, cmdfunc func(tgbotapi.Update)) {
	if c.commands == nil {
		c.commands = make(map[string]func(tgbotapi.Update))
	}
	c.commands[cmdstring] = cmdfunc
}

func (c *CommandHandler) RemoveCommand(cmdstring string) {
	delete(c.commands, cmdstring)
}

func (c *CommandHandler) GetCommandFunction(cmdstring string) func(tgbotapi.Update) {
	return c.commands[cmdstring]
}

func (c *CommandHandler) GetAndRunCommand(cmdstring string, update tgbotapi.Update) bool {
	fn := c.commands[cmdstring]
	if fn == nil {
		return false
	}
	fn(update)
	return true
}

func (bot *Bot) AddCommands() {
	// start command

	bot.CommandHandler.RegisterCommand("start", bot.SendStartMessage)

	// help command

	bot.CommandHandler.RegisterCommand("help", bot.SendHelpMessage)

	// hi command
	bot.CommandHandler.RegisterCommand("hi", bot.SendHi)
	//spamban command
	bot.CommandHandler.RegisterCommand("spamban", bot.HandleSpamMode)
	// ban
	bot.CommandHandler.RegisterCommand("ban", bot.BanUser)
	// unban
	bot.CommandHandler.RegisterCommand("unban", bot.UnbanUser)
	// kick
	bot.CommandHandler.RegisterCommand("kick", bot.KickUser)
	//testphoto
	bot.CommandHandler.RegisterCommand("testphoto", bot.SendTestPhoto)
	//filter
	bot.CommandHandler.RegisterCommand("filter", bot.AddFilter)
	// removefilter
	bot.CommandHandler.RegisterCommand("removefilter", bot.RemoveFilter)
	// listfilters
	bot.CommandHandler.RegisterCommand("filters", bot.AllFilters)
	// removeallfilters
	bot.CommandHandler.RegisterCommand("removeallfilters", bot.RemoveAllFilters)
	// addblacklist
	bot.CommandHandler.RegisterCommand("addblacklist", bot.AddBlackList)
	// addbottodbandroot
	bot.CommandHandler.RegisterCommand("addbottoroot", bot.AddBotToRoot)
	bot.CommandHandler.RegisterCommand("addbotindb", bot.AddBot)

	// Crypto Commands
	bot.CommandHandler.RegisterCommand("scan", bot.ScanContract)
}

func (bot *Bot) HandleCommand(update tgbotapi.Update) {
	message := update.Message
	if message.Command() == "" {
		return
	}
	if bot.CommandHandler.GetAndRunCommand(message.Command(), update) {
		return
	}
	// checlk and run filters!
	// if bot.CommandHandler.GetAndRunCommand(message.Command()+"@"+bot.Bot.Token, update) {
	// 	return
	// }

}

func (bot *Bot) HelpMessage(update tgbotapi.Update) {

}

func (bot *Bot) HandleSpamMode(update tgbotapi.Update) {
	if update.Message.Chat.Type == "private" {
		return
	}
}

func (bot *Bot) SendStartMessage(update tgbotapi.Update) {
	msg := bot.GetStartMessage(update.Message.Chat.ID, true)
	bot.Bot.Send(msg)
}

func (bot *Bot) SendHelpMessage(update tgbotapi.Update) {
	msg := bot.GetHelpMessage(update.Message.Chat.ID, true)
	bot.Bot.Send(msg)
}

func (bot *Bot) SendHi(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi :)")
	bot.Bot.Send(msg)
}
