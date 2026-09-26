package pubsub

import "time"

type mockedDriver struct {
	message *Message
	ticker  time.Ticker
	stop    chan struct{}
}

func NewMockedDriver(ticker time.Ticker, stop chan struct{}, mockedMessage *Message) DriverContract {
	return &mockedDriver{
		message: mockedMessage,
		ticker:  ticker,
		stop:    stop,
	}
}

func (m *mockedDriver) Push(topicName string, data []byte) (string, error) {
	return "", nil
}

func (m *mockedDriver) Subscribe(channel string, ch chan *Message) error {
	go func() {
		for {
			select {
			case <-m.ticker.C:
				ch <- m.message
			case <-m.stop:
				return
			}
		}

	}()
	return nil
}
