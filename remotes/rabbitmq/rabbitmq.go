package rabbitmq

import (
	"encoding/json"

	c "github.com/joaopandolfi/blackwhale/configurations"
	"github.com/streadway/amqp"
	"golang.org/x/xerrors"
)

// Driver for RabbitMQ
type Driver struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queues  map[string]*amqp.Queue
}

var conn *amqp.Connection
var chanel *amqp.Channel

func open() (*amqp.Connection, *amqp.Channel, error) {
	c, err := amqp.Dial(c.Configuration.RabbitMQURL)
	if err != nil {
		return nil, nil, xerrors.Errorf("connecting to rabbitmq:: %w", err)
	}

	ch, err := c.Channel()
	if err != nil {
		return nil, nil, xerrors.Errorf("opening a channel: %w", err)
	}

	return c, ch, nil
}

func new(c *amqp.Connection, ch *amqp.Channel) *Driver {
	return &Driver{
		Conn:    c,
		Channel: ch,
		Queues:  map[string]*amqp.Queue{},
	}
}

// New rabbitmq driver singleton
func New() (*Driver, error) {
	if conn == nil {
		c, ch, err := open()
		if err != nil {
			return nil, xerrors.Errorf("creating rabbit mq driver: %w", err)
		}
		conn = c
		chanel = ch
	}
	return new(conn, chanel), nil
}

// Fresh return a new fresh and clean connection
func Fresh() (*Driver, error) {
	c, ch, err := open()
	if err != nil {
		return nil, xerrors.Errorf("creating fresh rabbit mq driver: %w", err)
	}
	return new(c, ch), nil
}

func (d *Driver) OpenQueue(tube string) error {
	if d.Queues[tube] == nil {
		q, err := d.Channel.QueueDeclare(
			tube,  // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return xerrors.Errorf("declaring tube: %w", err)
		}
		d.Queues[tube] = &q
	}

	return nil
}

func (d *Driver) PutDefault(tube string, body interface{}) error {
	b, err := json.Marshal(&body)
	if err != nil {
		return xerrors.Errorf("marshaling body: %w", err)
	}
	err = d.Channel.Publish(
		"",    // exchange
		tube,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)

	if err != nil {
		return xerrors.Errorf("publishing on channel: %w", err)
	}

	return nil
}

func (d *Driver) Consume(tube string) (<-chan amqp.Delivery, error) {
	return d.Channel.Consume(
		tube,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
