package main

import (
	"SafeApeBot/bot"
	"SafeApeBot/config"
	"SafeApeBot/db"
	"SafeApeBot/server"
	"SafeApeBot/structs"
	"flag"
)

func Argshandler(conf *structs.SConfig) {
	serverport := flag.Int("port", conf.ServerPort, "Server Port")
	BotToken := flag.String("token", conf.BotToken, "Bot Token")
	flag.Parse()
	conf.ServerPort = *serverport
	conf.BotToken = *BotToken
}
func main() {
	// SetGlobals()
	config, err := config.Read_config_from_file("config.json")
	Argshandler(config)
	if err != nil {
		panic(err)
	}
	mdb := db.InitDb(config)
	_ = mdb
	// userrr := db.DbUserStruct{Name: "Test ACc", Userid: 123, SudoUser: true, Banned: false, IsPremium: true}
	// mdb.AddUser(userrr)
	// mdb.GetUserById(123)
	// mdb.TestDB()
	// fmt.Println(config)
	bot, err := bot.InitByToken(config.BotToken)
	if err != nil {
		panic(err)
	}
	bot.DB = &mdb
	bot.SetConfig(config)
	bot.SendText(1248191458, "Test")
	go bot.Start()

	// server
	server.DB = &mdb
	server.Config = config
	server := server.Server{}
	server.StartServer()
}
