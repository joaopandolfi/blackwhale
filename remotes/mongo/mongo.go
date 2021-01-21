package mongo

import (
	"github.com/joaopandolfi/blackwhale/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Index -
type Index struct {
	N   int    `json:"n"`
	Key string `json:"key"`
}

// GetSession - return session from a pool
func GetSession() *Session {
	session, err := GetPoolSession() //NewSession()
	//session, err := NewSession()
	if err != nil {
		utils.CriticalError("Unable to connect on mongo: %s", err)
		FlushPull()
		panic(err)
	}
	return session
}

// GenericInsert - insert new item on database
func GenericInsert(collection string, data interface{}) error {
	session := GetSession()
	return session.GetCollection(collection).Insert(&data)
}

// Run specific command
func Run(cmd interface{}) {
	session := GetSession()
	session.Run(cmd)
}

// CreateIndex create index on collection and key
func CreateIndex(collection string, keys ...string) error {
	session := GetSession()
	col := session.GetCollection(collection)
	col.EnsureIndexKey(keys...)
	return nil
}

// GetNextID returns next incremental id
func GetNextID(key string) (id int) {
	session := GetSession()
	col := session.GetCollection("whale_counter")

	var doc Index
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"n": 1}},
		ReturnNew: true,
	}
	_, err := col.Find(bson.M{"key": key}).Apply(change, &doc)
	if err != nil {
		err = col.Insert(bson.M{"key": key, "n": 0})
		if err != nil {
			utils.CriticalError("[Mongo][GetNextID] - Error on get Next ID", err)
			FlushPull()
			panic(err)
		}
		doc.N = 0
	}
	id = doc.N
	return
}

// Close all connections
func Close() {
	FlushPull()
}
