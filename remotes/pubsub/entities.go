package pubsub

import (
	ps "cloud.google.com/go/pubsub"
)

type Message struct {
	M *ps.Message
}

type subscription struct {
	channels []chan *Message
	sub      *ps.Subscription
}
