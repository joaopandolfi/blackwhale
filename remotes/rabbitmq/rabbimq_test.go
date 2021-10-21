package rabbitmq

import (
	"testing"

	c "github.com/joaopandolfi/blackwhale/configurations"
)

const testUrl = "amqp://guest:guest@10.0.0.2:5672/"

func Test_putdata(t *testing.T) {
	c.Configuration = c.Configurations{RabbitMQURL: testUrl}
	d, err := New()
	if err != nil {
		t.Errorf("creating rabbitMQ driver: %v", err)
		return
	}
	err = d.PutDefault("teste", map[string]string{"msg": "bananinha amassada"})
	if err != nil {
		t.Errorf("puting data on tube teste: %v", err)
		return
	}
}

func Test_readdata(t *testing.T) {
	c.Configuration = c.Configurations{RabbitMQURL: testUrl}
	d, err := New()
	if err != nil {
		t.Errorf("creating rabbitMQ driver: %v", err)
		return
	}

	c, err := d.Consume("teste")
	if err != nil {
		t.Errorf("consuming message from tube teste: %v", err)
		return
	}
	body := <-c
	t.Logf("%s, %s", body.AppId, string(body.Body))
}
