package structs

type (
	SConfig struct {
		Privatekey    string    `json:"privatekey"`
		Networks      []Network `json:"networks"`
		BotToken      string    `json:"BotToken"`
		BotGroup      string    `json:"BotGroup"`
		MongoDBURI    string    `json:"MongoDburl"`
		DatabaseToUse string    `json:"DB"`
		ServerPort    int       `json:"ServerPort"`
	}
	Network struct {
		NodeURL  string `json:"node_url"`
		ChainID  int    `json:"chain_id"`
		Name     string `json:"name"`
		Explorer string `json:"explorer"`
	}

	State map[string]map[int64]string
)

// Bot types
