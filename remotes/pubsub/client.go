package pubsub

// code imported from authenticator src/remotes/pubsub/ on version 0.1.4

import (
	"context"
	"fmt"
	"sync"

	ps "cloud.google.com/go/pubsub"
)

type Driver struct {
	Client        *ps.Client
	Ctx           context.Context
	subscriptions map[string]*subscription
	topics        map[string]*ps.Topic
	muTopics      sync.Mutex
	muSubs        sync.Mutex
}

var client *Driver

// Init start the pubsub client
// To work, you need to setup the env GOOGLE_APPLICATION_CREDENTIALS with the json path
// The json needs to contain the generated credentials on google cloud console
func Init(c context.Context, projectID string) error {
	cl, err := ps.NewClient(c, projectID)
	if err != nil {
		return fmt.Errorf("connecting to PubSub: %w", err)
	}

	client = &Driver{
		Client:        cl,
		Ctx:           c,
		subscriptions: map[string]*subscription{},
		topics:        map[string]*ps.Topic{},
	}
	return nil
}

func Get() *Driver {
	if client == nil {
		panic("PubSub client driver not initialized")
	}
	return client
}

func Close() {
	if client != nil && client.Client != nil {
		client.Client.Close()
	}
}

func (c *Driver) getTopic(topic string) *ps.Topic {
	c.muTopics.Lock()
	defer c.muTopics.Unlock()
	if topic, ok := c.topics[topic]; ok {
		return topic
	}
	t := c.Client.Topic(topic)
	c.topics[topic] = t
	return t
}

// Push send message to pubsub
// returns messageID and error
func (c *Driver) Push(topicName string, data []byte) (string, error) {
	topic := c.getTopic(topicName)
	res := topic.Publish(c.Ctx, &ps.Message{Data: data})
	msgID, err := res.Get(c.Ctx)
	if err != nil {
		return "", fmt.Errorf("publishing on topic %s: %w", topicName, err)
	}
	return msgID, nil
}

func (c *Driver) Subscribe(channel string, ch chan *Message) error {
	c.muSubs.Lock()
	defer c.muSubs.Unlock()
	if c.Client == nil {
		return fmt.Errorf("client is not inialized")
	}

	if _channel, ok := c.subscriptions[channel]; ok {
		_channel.channels = append(_channel.channels, ch)
		return nil
	}
	sub := c.Client.Subscription(channel)
	c.subscriptions[channel] = &subscription{
		sub:      sub,
		channels: []chan *Message{ch},
	}

	err := sub.Receive(c.Ctx, func(ctx context.Context, m *ps.Message) {
		for i := range c.subscriptions[channel].channels {
			c.subscriptions[channel].channels[i] <- &Message{M: m}
		}
	})

	if err != nil {
		return fmt.Errorf("injecting receiver on subscription: %w", err)
	}

	return nil
}

/* === REFFERENCE
// Publish "hello world" on topic1.
topic := client.Topic("topic1")
res := topic.Publish(ctx, &pubsub.Message{
	Data: []byte("hello world"),
})
// The publish happens asynchronously.
// Later, you can get the result from res:
...
msgID, err := res.Get(ctx)
if err != nil {
	log.Fatal(err)
}

// Use a callback to receive messages via subscription1.
sub := client.Subscription("subscription1")
err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
	fmt.Println(m.Data)
	m.Ack() // Acknowledge that we've consumed the message.
})
if err != nil {
	log.Println(err)
}
*/
