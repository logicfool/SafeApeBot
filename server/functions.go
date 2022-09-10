package server

import (
	"SafeApeBot/server/bothook"
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/julienschmidt/httprouter"
)

// also pass a state var in the bot to be able to see if the user is in a state of typing something or to handle stuff like that!
func HandleBotWebHook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// res := "Bot Token: " + ps.ByName("token")
	// fmt.Println(res)
	// fmt.Fprintf(w, res)
	// return
	var update tgbotapi.Update
	var bot bothook.Bot
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		fmt.Println("Error :", err)
	}
	// spew.Dump(update.Message)
	bot, err = bothook.InitFromTokenAndUpdate(ps.ByName("token"), update, w)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	bot.BotSession.DB = DB
	if bot.BotSession.Config == nil {
		bot.BotSession.SetConfig(Config)
		bot.Config = Config
	}
	bot.HandleUpdate()
	fmt.Fprintf(w, "200")
	// fmt.Fprintf(w, "Hello! Bottest")
}
