package v2

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client struct {
	conn *mongo.Client
}

func MountURL(username, password, host, options string) string {
	if username == "" || password == "" {
		return fmt.Sprintf("mongodb://%s/%s", host, options)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, host, options)
}

// New mongo client
func New(url string) (*Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo: %w", err)
	}

	return &Client{conn: client}, nil
}

// Disconnect client to server
func (c *Client) Disconnect(ctx context.Context) error {
	if c.conn != nil {
		return c.conn.Disconnect(ctx)
	}

	return nil
}

// Conn - client connection
func (c *Client) Conn() (*mongo.Client, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("database is not setted")
	}
	return c.conn, nil
}

func (c *Client) Collection(database, collection string) *mongo.Collection {
	return c.conn.Database(database).Collection(collection)
}

type index struct {
	N   int    `json:"n"`
	Key string `json:"key"`
}

// GetNextCounter returns next incremental counter
func (c *Client) GetNextCounter(ctx context.Context, database, key string) (int, error) {
	var doc index
	coll := c.Collection(database, "whale_counter")

	filter := bson.M{"key": key}
	update := bson.M{"$inc": bson.M{"n": 1}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		_, err = coll.InsertOne(ctx, bson.M{"key": key, "n": 0})
		if err != nil {
			return 0, fmt.Errorf("initializing counter (%s): %w", key, err)
		}
		doc.N = 0
	}
	return doc.N, nil
}

// ClearCounter reset counter
func (c *Client) ClearCounter(ctx context.Context, database, key string) error {
	coll := c.Collection(database, "whale_counter")
	_, err := coll.DeleteOne(ctx, bson.M{"key": key})
	if err != nil {
		return fmt.Errorf("deleting counter: %w", err)
	}
	return nil
}
