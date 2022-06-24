package utils

import (
	"SafeApeBot/db"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type (
	CacheStruct struct {
		Cache *cache.Cache
	}
)

var Cache CacheStruct = CacheStruct{Cache: cache.New(5*time.Minute, 10*time.Minute)}

func (cache *CacheStruct) GetUser(userid int64) db.DBUserStruct {
	cuser := "user:" + strconv.Itoa(int(userid))
	res, success := cache.Cache.Get(cuser)
	if success {
		return res.(db.DBUserStruct)
	}
	return db.DBUserStruct{}
}

func (cache *CacheStruct) GetBot(btoken string) db.DBBotStruct {
	cchat := "bot:" + btoken
	res, success := cache.Cache.Get(cchat)
	if success {
		return res.(db.DBBotStruct)
	}
	return db.DBBotStruct{}
}

func (cache1 *CacheStruct) SetUser(userid int64, user db.DBUserStruct) {
	cuser := "user:" + strconv.Itoa(int(userid))
	cache1.Cache.Set(cuser, user, cache.NoExpiration)
}

func (cache1 *CacheStruct) SetBot(btoken string, chat db.DBBotStruct) {
	cchat := "bot:" + btoken
	cache1.Cache.Set(cchat, chat, cache.NoExpiration)
}
