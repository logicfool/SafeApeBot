package bot

import (
	"SafeApeBot/db"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) BanUser(update tgbotapi.Update) {
	// from_user := update.Message.From
	chat_id := update.Message.Chat.ID
	isadmin := bot.IsChatAdmin(update)
	if !isadmin {
		bot.SendText(chat_id, "You are not admin")
		return
	}
	reply_to_message := update.Message.ReplyToMessage
	var user_to_ban int64
	if reply_to_message != nil {
		user_to_ban = reply_to_message.From.ID
		userconf := tgbotapi.ChatMemberConfig{ChatID: chat_id, UserID: user_to_ban}
		ban := tgbotapi.BanChatMemberConfig{ChatMemberConfig: userconf}
		bot.Bot.Send(ban)
		msg := tgbotapi.NewMessage(chat_id, fmt.Sprintf("%s has been banned", reply_to_message.From.FirstName))
		bot.Bot.Send(msg)
		return
	} else {
		return
	}

}

func (bot *Bot) KickUser(update tgbotapi.Update) {
	// from_user := update.Message.From
	chat_id := update.Message.Chat.ID
	isadmin := bot.IsChatAdmin(update)
	if !isadmin {
		bot.SendText(chat_id, "You are not admin")
		return
	}
	reply_to_message := update.Message.ReplyToMessage
	var user_to_ban int64
	if reply_to_message != nil {
		user_to_ban = reply_to_message.From.ID
		userconf := tgbotapi.ChatMemberConfig{ChatID: chat_id, UserID: user_to_ban}
		timeinepoch := int64(time.Now().Unix()) + 40
		ban := tgbotapi.BanChatMemberConfig{ChatMemberConfig: userconf, UntilDate: timeinepoch}
		bot.Bot.Send(ban)
		msg := tgbotapi.NewMessage(chat_id, fmt.Sprintf("%s has been kicked!", reply_to_message.From.FirstName))
		bot.Bot.Send(msg)
		return
	}
}
func (bot *Bot) UnbanUser(update tgbotapi.Update) {
	// from_user := update.Message.From
	chat_id := update.Message.Chat.ID
	isadmin := bot.IsChatAdmin(update)
	if !isadmin {
		bot.SendText(chat_id, "You are not admin")
		return
	}
	reply_to_message := update.Message.ReplyToMessage
	var user_to_ban int64
	if reply_to_message != nil {
		user_to_ban = reply_to_message.From.ID
		userconf := tgbotapi.ChatMemberConfig{ChatID: chat_id, UserID: user_to_ban}
		unban := tgbotapi.UnbanChatMemberConfig{ChatMemberConfig: userconf, OnlyIfBanned: true}
		bot.Bot.Send(unban)
		msg := tgbotapi.NewMessage(chat_id, fmt.Sprintf("%s has been unbanned", reply_to_message.From.FirstName))
		bot.Bot.Send(msg)
		return
	}
}

func (bot *Bot) SendTestPhoto(update tgbotapi.Update) {
	chat_id := update.Message.Chat.ID
	// photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg"))
	// mediaGroup := tgbotapi.NewMediaGroup(chat_id, []interface{}{
	// 	photo,
	// })
	photo := tgbotapi.NewPhoto(chat_id, tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg"))
	photo.Caption = "This is a test photo"
	bot.Bot.Send(photo)
}

func (bot *Bot) AddFilter(update tgbotapi.Update) {
	slices := strings.Split(update.Message.Text, " ")
	bot_in_db := bot.DB.GetBotByToken(bot.Bot.Token)
	spew.Dump(bot_in_db)
	chats := bot_in_db.Chats
	if chats == nil {
		chats := make(map[int64]db.DBChatStruct)
		bot_in_db.Chats = chats
	}
	chat := chats[update.Message.Chat.ID]
	if bot_in_db.BotToken == "" {
		return
	}
	spew.Dump(bot_in_db)
	if chat.ChatID == 0 {
		chat.ChatID = update.Message.Chat.ID
		chat.IsActive = true
		// bot.DB.UpdateBot(bot_in_db)
		// bot_in_db = bot.DB.GetBotByToken(bot.Bot.Token)
		// chats = bot_in_db.Chats
		// chat = chats[update.Message.Chat.ID]
		// spew.Dump(chat)
	}
	if chat.Filters == nil {
		chat.Filters = make(map[string]string)
		chats[update.Message.Chat.ID] = chat
		bot_in_db.Chats = chats
	}
	if len(slices) == 2 {
		filter_word := strings.ToLower(slices[1])
		if update.Message.ReplyToMessage != nil {
			var filter_text string
			// might be a gif or a video or photo or text handle it based on that
			if update.Message.ReplyToMessage.Photo != nil {
				filter_text = "photo::" + update.Message.ReplyToMessage.Photo[0].FileID + "::" + update.Message.ReplyToMessage.Caption
			} else if update.Message.ReplyToMessage.Animation != nil {
				filter_text = "animation::" + update.Message.ReplyToMessage.Animation.FileID + "::" + update.Message.ReplyToMessage.Caption
			} else if update.Message.ReplyToMessage.Video != nil {
				filter_text = "video::" + update.Message.ReplyToMessage.Video.FileID + "::" + update.Message.ReplyToMessage.Caption
			} else if update.Message.ReplyToMessage.Audio != nil {
				filter_text = "audio::" + update.Message.ReplyToMessage.Audio.FileID + "::" + update.Message.ReplyToMessage.Caption
			} else if update.Message.ReplyToMessage.Document != nil {
				filter_text = "document::" + update.Message.ReplyToMessage.Document.FileID + "::" + update.Message.ReplyToMessage.Caption
			} else if update.Message.ReplyToMessage.Sticker != nil {
				filter_text = "sticker::" + update.Message.ReplyToMessage.Sticker.FileID
			} else {
				filter_text = "text::" + update.Message.ReplyToMessage.Text
			}
			bot_in_db.Chats[update.Message.Chat.ID].Filters[filter_word] = filter_text
			bot.DB.UpdateBot(bot_in_db)
			bot.SendText(update.Message.Chat.ID, "Filter "+filter_word+" added")
			return
		}
	} else if len(slices) > 2 {
		filter_word := strings.ToLower(slices[1])
		// remove 0 and 1 index from slices
		slices = slices[2:]
		filter_text := strings.Join(slices, " ")
		bot_in_db.Chats[update.Message.Chat.ID].Filters[filter_word] = "text::" + filter_text
		spew.Dump(bot_in_db)
		bot.DB.UpdateBot(bot_in_db)
		bot.SendText(update.Message.Chat.ID, "Filter "+filter_word+" added")
		return
	}
}

func (bot *Bot) RemoveFilter(update tgbotapi.Update) {
	slices := strings.Split(update.Message.Text, " ")
	bot_in_db := bot.DB.GetBotByToken(bot.Bot.Token)
	if bot_in_db.BotToken == "" {
		return
	}
	chats := bot_in_db.Chats
	if chats == nil {
		chats := make(map[int64]db.DBChatStruct)
		bot_in_db.Chats = chats
	}
	chat := chats[update.Message.Chat.ID]
	if chat.ChatID == 0 {
		chat.ChatID = update.Message.Chat.ID
		chat.IsActive = true
		chats[update.Message.Chat.ID] = chat
		bot_in_db.Chats = chats
		bot.DB.UpdateBot(bot_in_db)
		return
	}
	if chat.Filters == nil {
		return
	}
	if len(slices) == 2 {
		filter_word := slices[1]
		filterr := chat.Filters[filter_word]
		if filterr == "" {
			return
		}
		delete(chat.Filters, filter_word)
		chats[update.Message.Chat.ID] = chat
		bot_in_db.Chats = chats
		bot.DB.UpdateBot(bot_in_db)
		bot.SendText(update.Message.Chat.ID, "Filter removed")
		return
	}
}

func (bot *Bot) AllFilters(update tgbotapi.Update) {
	bot_in_db := bot.DB.GetBotByToken(bot.Bot.Token)
	if bot_in_db.BotToken == "" {
		// fmt.Println("here!")
		return
	}
	chats := bot_in_db.Chats
	if chats == nil {
		chats := make(map[int64]db.DBChatStruct)
		bot_in_db.Chats = chats
	}
	chat := chats[update.Message.Chat.ID]
	if chat.ChatID == 0 {
		chat.ChatID = update.Message.Chat.ID
		chat.IsActive = true
		chats[update.Message.Chat.ID] = chat
		bot_in_db.Chats = chats
		bot.DB.UpdateBot(bot_in_db)
		bot.SendText(update.Message.Chat.ID, "No Filters found in this chat!")
		return
	}
	if chat.Filters == nil {
		bot.SendText(update.Message.Chat.ID, "No Filters found in this chat!")
		return
	}
	var text string = "Filters in this chat:\n"
	for filter := range chat.Filters {
		text += "`" + filter + "`" + "\n"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Bot.Send(msg)
}

func (bot *Bot) SearchAdnSendFilter(update tgbotapi.Update) {
	text := strings.ToLower(update.Message.Text)
	bot_in_db := bot.DB.GetBotByToken(bot.Bot.Token)
	if bot_in_db.BotToken == "" {
		return
	}
	chats := bot_in_db.Chats
	if chats == nil {
		chats := make(map[int64]db.DBChatStruct)
		bot_in_db.Chats = chats
	}
	chat := chats[update.Message.Chat.ID]
	if chat.ChatID == 0 {
		chat.ChatID = update.Message.Chat.ID
		chat.IsActive = true
		chats[update.Message.Chat.ID] = chat
		bot_in_db.Chats = chats
		bot.DB.UpdateBot(bot_in_db)
		return
	}
	if chat.Filters == nil {
		return
	}
	slices_of_msg := strings.Split(text, " ")

	for _, msg_word := range slices_of_msg {
		// msg_word = strings.ToLower(msg_word)
		filter := chat.Filters[msg_word]
		if filter == "" {
			continue
		}
		slices := strings.Split(filter, "::")
		if slices[0] == "text" {
			// bot.SendText(update.Message.Chat.ID, slices[1])
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, slices[1])
			msg.ParseMode = "Markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Bot.Send(msg)
		} else if slices[0] == "photo" {
			photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				photo.Caption = slices[2]
			}
			photo.ReplyToMessageID = update.Message.MessageID
			bot.Bot.Send(photo)
		} else if slices[0] == "audio" {
			audio := tgbotapi.NewAudio(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				audio.Caption = slices[2]
			}
			audio.ReplyToMessageID = update.Message.MessageID
		} else if slices[0] == "video" {
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				video.Caption = slices[2]
			}
			video.ReplyToMessageID = update.Message.MessageID
		} else if slices[0] == "document" {
			document := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				document.Caption = slices[2]
			}
			document.ReplyToMessageID = update.Message.MessageID
		} else if slices[0] == "sticker" {
			sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			sticker.ReplyToMessageID = update.Message.MessageID
		} else if slices[0] == "voice" {
			voice := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				voice.Caption = slices[2]
			}
			voice.ReplyToMessageID = update.Message.MessageID
		} else if slices[0] == "animation" {
			animation := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileID(slices[1]))
			if slices[2] != "" {
				animation.Caption = slices[2]
			}
			animation.ReplyToMessageID = update.Message.MessageID
			bot.Bot.Send(animation)
		}
		return
	}

}
