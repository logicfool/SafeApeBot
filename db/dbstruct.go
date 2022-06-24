package db

type (
	DBUserStruct struct {
		// ID        string `bson:"_id"`
		Name          string `json:"name"`
		Userid        int64  `json:"userid"`
		SudoUser      bool   `json:"is_sudo"`
		Banned        bool   `json:"banned"`
		IsPremium     bool   `json:"is_premium"`
		HasStartedBot bool   `json:"bot_started"`
		Currstate     string
	}
	DBBotStruct struct {
		// ID       string                 `bson:"_id"`
		BotToken string                 `json:"bottoken"`
		SudoUser int64                  `json:"created_by"`
		Chats    map[int64]DBChatStruct `json:"chats"`
		RootBot  bool                   `json:"sudo_bot"`
		Enabled  bool                   `json:"enabled"`
	}

	DBChatStruct struct {
		// ID                   string                 `bson:"_id"`
		ChatID               int64                  `json:"chat_id"`
		Filters              map[string]string      `json:"filters"`
		Blacklist            map[string]interface{} `json:"blacklist"`
		KickSpamMode         bool                   `json:"spam_mode"`
		Warnings             map[int64]int          `json:"warnings"`
		VerificationType     int                    `json:"verification_mode"`
		ETHAddressAutoRemove bool                   `json:"eth_address_auto_remove"`
		IsActive             bool                   `json:"is_active"`
	}

	DBMempoolStruct struct {
		ID              string `bson:"_id"`
		Network         string `json:"network"`
		Address         string `json:"address"`
		Type            string `json:"type"`
		AddedbyUser     int    `json:"from_user"`
		ChatID          int    `json:"chat_id"`
		BotToken        string `json:"bottoken"`
		ShowSells       bool   `json:"show_sells"`
		MessageTemplate string `json:"message_template"`
	}
)
