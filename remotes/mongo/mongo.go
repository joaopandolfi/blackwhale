package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/joaopandolfi/blackwhale/v2/utils"
)

type counterDoc struct {
	N   int    `bson:"n" json:"n"`
	Key string `bson:"key" json:"key"`
}

// GetSession - return the shared session for the configured mongo url
func GetSession() *Session {
	session, err := GetPoolSession()
	if err != nil {
		utils.CriticalError("Unable to connect on mongo: %s", err)
		Close()
		panic(err)
	}
	return session
}

// CreateIndex create a compound ascending index on collection
func CreateIndex(collection string, keys ...string) error {
	session := GetSession()
	col := session.GetCollection(collection)

	indexKeys := bson.D{}
	for _, k := range keys {
		indexKeys = append(indexKeys, bson.E{Key: k, Value: 1})
	}
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: indexKeys})
	return err
}

// GetNextID returns next incremental id
func GetNextID(key string) (id int) {
	session := GetSession()
	col := session.GetCollection("whale_counter")
	ctx := context.Background()

	var doc counterDoc
	err := col.FindOneAndUpdate(ctx,
		bson.M{"key": key},
		bson.M{"$inc": bson.M{"n": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		_, iErr := col.InsertOne(ctx, bson.M{"key": key, "n": 0})
		if iErr != nil {
			utils.CriticalError("[Mongo][GetNextID] - Error on get Next ID", iErr)
			Close()
			panic(iErr)
		}
		return 0
	}
	return doc.N
}

// Close all connections
func Close() {
	mu.Lock()
	defer mu.Unlock()
	for url, client := range clients {
		client.Disconnect(context.Background())
		delete(clients, url)
	}
}
