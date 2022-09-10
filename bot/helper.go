package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) IsChatAdmin(update tgbotapi.Update) bool {
	if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
		chatconfigwithuser := tgbotapi.ChatConfigWithUser{ChatID: update.Message.Chat.ID, UserID: update.Message.From.ID}
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

func (bot *Bot) IsMember(update tgbotapi.Update) bool {
	// update := update
	if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
		chatconfigwithuser := tgbotapi.ChatConfigWithUser{ChatID: update.Message.Chat.ID, UserID: update.Message.From.ID}
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

func (bot *Bot) handleNewMemberJoin(update tgbotapi.Update) {
	newchatmembers := update.Message.NewChatMembers
	_ = newchatmembers
}

func (bot *Bot) IsServiceMessage(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
	if update.Message.ReplyToMessage != nil || update.Message.Text != "" || update.Message.Animation != nil || update.Message.Audio != nil || update.Message.Document != nil || update.Message.Photo != nil || update.Message.Sticker != nil || update.Message.Video != nil || update.Message.Voice != nil || update.Message.VideoNote != nil || update.Message.Contact != nil || update.Message.Location != nil || update.Message.Venue != nil || update.Message.Caption != "" || update.Message.Dice != nil || update.Message.Game != nil || update.Message.Poll != nil {
		return false
	}

	return true
}

func (bot *Bot) GetChatByusername(username string) tgbotapi.Chat {
	// username = strings.Replace(username, "@", "", -1)
	// fmt.Println(userrname)
	if !strings.Contains(username, "@") {
		username = "@" + username
	}
	user := tgbotapi.ChatConfig{SuperGroupUsername: username}
	chatinfoconfig := tgbotapi.ChatInfoConfig{ChatConfig: user}
	chat, err := bot.Bot.GetChat(chatinfoconfig)
	if err != nil {
		log.Println(err)
	}
	return chat
}
