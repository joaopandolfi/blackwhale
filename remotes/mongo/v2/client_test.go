//go:build integration

package v2

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var _host = envOr("MONGO_HOST", "127.0.0.1")
var _port = envOr("MONGO_PORT", "27017")
var _database = envOr("MONGO_DATABASE", "black_whale_test")

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Test_Connection(t *testing.T) {
	url := MountURL("", "", _host, _port)
	client, err := New(url, nil)
	if err != nil {
		t.Errorf("connecting: %v", err)
		return
	}

	if err := client.conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Errorf("ping: %v", err)
		return
	}
	err = client.Disconnect()
	if err != nil {
		t.Errorf("disconnect: %v", err)
		return
	}
}

func Test_Manipulating(t *testing.T) {
	url := MountURL("", "", _host, _port)
	client, err := New(url, nil)
	if err != nil {
		t.Errorf("connecting: %v", err)
		return
	}

	col := client.Collection(_database, "test")

	data := bson.M{"title": "test", "text": "trash"}

	_, err = col.InsertOne(context.TODO(), data)
	if err != nil {
		t.Errorf("inserting: %v", err)
		return
	}

	filter := bson.M{"title": "test"}
	update := bson.M{"$set": bson.M{"title": "updated"}}
	_, err = col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		t.Errorf("updating: %v", err)
		return
	}

	var result map[string]string

	err = col.FindOne(context.TODO(), bson.M{"title": "updated"}).Decode(&result)
	if err != nil {
		t.Errorf("searching: %v", err)
		return
	}

	if result["title"] != "updated" {
		t.Errorf("title different from expected: (given: %s expected: %s)", result["title"], "updated")
		return
	}

	filter = bson.M{"title": "updated"}
	_, err = col.DeleteOne(context.TODO(), filter)
	if err != nil {
		t.Errorf("deleting: %v", err)
		return
	}

	err = client.Disconnect()
	if err != nil {
		t.Errorf("disconnect: %v", err)
		return
	}
}

func Test_Counter(t *testing.T) {
	url := MountURL("", "", _host, _port)
	client, err := New(url, nil)
	if err != nil {
		t.Errorf("connecting: %v", err)
		return
	}

	_counter_key := "test"

	count, err := client.GetNextCounter(_database, _counter_key)
	if err != nil {
		t.Errorf("first count err: %v", err)
		return
	}

	if count != 0 {
		t.Errorf("expected: 0 got: %v", count)
		return
	}

	count, err = client.GetNextCounter(_database, _counter_key)
	if err != nil {
		t.Errorf("second count err: %v", err)
		return
	}

	if count != 1 {
		t.Errorf("expected: 1 got: %v", count)
		return
	}

	err = client.ClearCounter(_database, _counter_key)
	if err != nil {
		t.Errorf("cleaning counter: %v", err)
		return
	}
}
