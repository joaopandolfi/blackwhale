package mongo

import (
	"gopkg.in/mgo.v2"
)

// Collection -
type Collection struct {
	collection *mgo.Collection
}

// NewService -  create a collection and generate index
func NewService(session *Session, collectionName string, index mgo.Index) *Collection {
	collection := session.GetCollection(collectionName)
	collection.EnsureIndex(index)
	return &Collection{collection}
}
