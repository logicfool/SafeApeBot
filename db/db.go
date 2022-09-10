package db

import (
	"SafeApeBot/structs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	DB struct {
		db     *mongo.Client
		config *structs.SConfig
	}
)

func InitDb(conf *structs.SConfig) DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoDBURI))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return DB{db: client, config: conf}
}

func (db *DB) GetObjectIdfilter(object_id string) primitive.M {
	objectID, _ := primitive.ObjectIDFromHex(object_id)
	filter := bson.M{"_id": objectID}
	return filter
}

func (db DB) GetDb() *mongo.Client {
	return db.db
}

func (db *DB) GetUserById(userid int64) DBUserStruct {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userid", Value: userid}}
	coll, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	for coll.Next(ctx) {
		var result DBUserStruct
		err = coll.Decode(&result)
		if err != nil {
			panic(err)
		}
		return result
	}
	return DBUserStruct{}
}

func (db *DB) AddUser(user DBUserStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (db *DB) UpdateUser(user DBUserStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userid", Value: user.Userid}}
	update := bson.D{{Key: "$set", Value: user}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (db *DB) GetBotByToken(BotToken string) DBBotStruct {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("bots")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "bottoken", Value: BotToken}}
	coll, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	for coll.Next(ctx) {
		var result DBBotStruct
		err = coll.Decode(&result)
		if err != nil {
			panic(err)
		}
		return result
	}
	return DBBotStruct{}
}

func (db *DB) AddBot(bot DBBotStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("bots")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, bot)
	return err
}

func (db *DB) UpdateBot(bot DBBotStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("bots")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "bottoken", Value: bot.BotToken}}
	update := bson.D{{Key: "$set", Value: bot}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (db *DB) FindMempoolByUser(user int64) []DBMempoolStruct {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("mempool")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userid", Value: user}}
	coll, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var result []DBMempoolStruct
	for coll.Next(ctx) {
		var result1 DBMempoolStruct
		err = coll.Decode(&result1)
		if err != nil {
			panic(err)
		}
		result = append(result, result1)
	}
	result = nil
	return result
}

func (db *DB) FindMempoolByChat(chat int64) []DBMempoolStruct {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("mempool")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "chatd_id", Value: chat}}
	coll, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var result []DBMempoolStruct
	for coll.Next(ctx) {
		var result1 DBMempoolStruct
		err = coll.Decode(&result1)
		if err != nil {
			panic(err)
		}
		result = append(result, result1)
	}
	result = nil
	return result
}

func (db *DB) AddMempool(mempool DBMempoolStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("mempool")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, mempool)
	return err
}

func (db *DB) UpdateMempool(mempool DBMempoolStruct) error {
	collection := db.db.Database(db.config.DatabaseToUse).Collection("mempool")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// filter := bson.D{{Key: "chatd_id", Value: mempool.ChatdId}}
	filter := db.GetObjectIdfilter(mempool.ID)
	update := bson.D{{Key: "$set", Value: mempool}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
