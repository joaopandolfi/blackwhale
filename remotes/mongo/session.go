package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/joaopandolfi/blackwhale/v2/configurations"
)

// Session wraps a mongo-driver client
type Session struct {
	client *mongo.Client
}

var (
	clients map[string]*mongo.Client
	mu      sync.RWMutex
)

func connect(mongoURL string) (*mongo.Client, error) {
	url := strings.ReplaceAll(mongoURL, "ssl=true", "tls=true")
	opts := options.Client().ApplyURI(url)
	if strings.Contains(url, "tls=true") {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo: %w", err)
	}
	return client, nil
}

// GetPoolSession - return the shared session for the configured mongo url
func GetPoolSession() (*Session, error) {
	return GetCustomPoolSession(configurations.Configuration.MongoUrl)
}

// GetCustomPoolSession - return the shared session for mongoURL
func GetCustomPoolSession(mongoURL string) (*Session, error) {
	mu.RLock()
	client, ok := clients[mongoURL]
	mu.RUnlock()
	if ok && client != nil {
		return &Session{client: client}, nil
	}

	mu.Lock()
	defer mu.Unlock()
	if client, ok := clients[mongoURL]; ok && client != nil {
		return &Session{client: client}, nil
	}

	client, err := connect(mongoURL)
	if err != nil {
		return nil, err
	}
	clients[mongoURL] = client
	return &Session{client: client}, nil
}

// NewSession - return the shared session for the configured mongo url
func NewSession() (*Session, error) {
	return GetCustomPoolSession(configurations.Configuration.MongoUrl)
}

// NewCustomSessionFresh - create a fresh session (not shared) for mongoURL
func NewCustomSessionFresh(mongoURL string) (*Session, error) {
	client, err := connect(mongoURL)
	if err != nil {
		return nil, err
	}
	return &Session{client: client}, nil
}

func (s *Session) Copy() *Session {
	return &Session{client: s.client}
}

func (s *Session) GetCollection(col string) *mongo.Collection {
	return s.client.Database(configurations.Configuration.MongoDb).Collection(col)
}

func (s *Session) GetCollectionOnDB(db, col string) *mongo.Collection {
	return s.client.Database(db).Collection(col)
}

func (s *Session) Run(cmd any) {
	var result bson.M
	s.client.Database(configurations.Configuration.MongoDb).RunCommand(context.Background(), toOrderedCommand(cmd)).Decode(&result)
	fmt.Println(result)
}

func toOrderedCommand(cmd any) any {
	if m, ok := cmd.(map[string]any); ok {
		d := make(bson.D, 0, len(m))
		for k, v := range m {
			d = append(d, bson.E{Key: k, Value: v})
		}
		return d
	}
	return cmd
}

func (s *Session) Close() {
	if s.client != nil {
		s.client.Disconnect(context.Background())
	}
}

func (s *Session) Health() error {
	if s.client == nil {
		return fmt.Errorf("checking health: session is nil")
	}
	return s.client.Ping(context.Background(), nil)
}
