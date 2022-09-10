package bothook

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) IsChatAdmin() bool {
	// update := bot.Update
	if bot.Update.Message.Chat.Type == "group" || bot.Update.Message.Chat.Type == "supergroup" {
		chatconfigwithuser := tgbotapi.ChatConfigWithUser{ChatID: bot.Update.Message.Chat.ID, UserID: bot.Update.Message.From.ID}
		chatmemberconfig := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: chatconfigwithuser}
		member, err := bot.Bot.GetChatMember(chatmemberconfig)
		if err != nil {
			log.Println(err)
		}
		if member.Status == "creator" || member.Status == "administrator" {
			return true
		}
		return false
	}
	return false
}

func (bot *Bot) IsMember() bool {
	// update := bot.Update
	if bot.Update.Message.Chat.Type == "group" || bot.Update.Message.Chat.Type == "supergroup" {
		chatconfigwithuser := tgbotapi.ChatConfigWithUser{ChatID: bot.Update.Message.Chat.ID, UserID: bot.Update.Message.From.ID}
		chatmemberconfig := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: chatconfigwithuser}
		member, err := bot.Bot.GetChatMember(chatmemberconfig)
		if err != nil {
			log.Println(err)
		}
		if member.Status != "kicked" && member.Status != "left" {
			return true
		}
		return false
	}
	return false
}
