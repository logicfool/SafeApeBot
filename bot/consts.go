package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	StartMessage string = "Hello, I am %name% .\nI am a Multi-Purpose Bot.\nSome of my features:\n	- I can help you to manage your group.\n	-I can track wallets as well as get buy/sell txs of tokens!\n	- Can get the price of a token\n	-Do Honeypot/Extra Checks!\nAnd much more check the help section for more!\nIf you face any issues you can always join @safeapechat\nI hope you like me!\nThanks ❤️"
	HelpMessage  string = "<b>HELP</b>\n\nHello, I am %name%.\nI am a Group manager with alot of extra features to help you manage your groups as well as protect them. \n<b>To manage a group add the bot in the group and give it proper admin rights.</b>\n\nAll Available Commands of the bot:\n\n<b>Group Administration Commands:</b>\n/kick,kickme,/ban,/filter,/filters,/mute,/muteall,/setwelcome,/welcome (on/off),/pin,/unpin,/joinverify (on/off),/preverify\n\n<b>Defi Tools: </b>\n/hp (tokenaddress) : Do a Honeypot check on the token\n/price (tokenaddress) : Get the price of a token\n/settoken (tokenaddress) : Enable notification for buy/sell of token!\n/removetoken (tokenaddress) : Remove and Disable Notification of buy/sell of token\n/setwatchlist (address) : Watch Wallet address for txs and get notified\n/removewatchlist (address) : Remove a wallet from watchlist\n\n"
)

func (bot *Bot) GetStartMessage(chatid int64, rootbot bool) tgbotapi.MessageConfig {
	if bot.Me == (tgbotapi.User{}) {
		bot.GetMe()
	}
	msg := tgbotapi.NewMessage(chatid, "")
	// msg.Text = StartMessage
	msg.Text = strings.Replace(StartMessage, "%name%", bot.Me.FirstName+" "+bot.Me.LastName, -1)
	// msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ParseMode = tgbotapi.ModeHTML
	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	var keymarkup tgbotapi.InlineKeyboardMarkup
	button1 := tgbotapi.NewInlineKeyboardButtonData("Help👁‍🗨", "action:help")
	button2 := tgbotapi.NewInlineKeyboardButtonURL("Add to your group", "http://t.me/"+bot.Me.UserName+"?startgroup=true")
	row1 = append(row1, button1)
	row2 = append(row2, button2)
	button3 := tgbotapi.NewInlineKeyboardButtonURL("Join SafeApe Chat💬", "https://t.me/"+bot.Config.BotGroup)
	row2 = append(row2, button3)
	if rootbot {
		keymarkup = tgbotapi.NewInlineKeyboardMarkup(row1, row2)
	} else {
		keymarkup = tgbotapi.NewInlineKeyboardMarkup(row1)
	}
	msg.ReplyMarkup = &keymarkup
	return msg
}

func (bot *Bot) GetHelpMessage(chatid int64, rootbot bool) tgbotapi.MessageConfig {
	if bot.Me == (tgbotapi.User{}) {
		bot.GetMe()
	}
	msg := tgbotapi.NewMessage(chatid, "")
	msg.Text = strings.Replace(HelpMessage, "%name%", bot.Me.FirstName+" "+bot.Me.LastName, -1)
	if rootbot {
		msg.Text = msg.Text + "<b>Clone: </b>\n/clone : Get your own bot with all the features and a name you like!\n\n"
	}
	msg.Text = msg.Text + "Thanks!"
	// msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}
