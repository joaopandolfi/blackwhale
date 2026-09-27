//go:build integration

package rabbitmq

import (
	"os"
	"testing"

	c "github.com/joaopandolfi/blackwhale/v2/configurations"
)

func testURL() string {
	if v := os.Getenv("RABBITMQ_URL"); v != "" {
		return v
	}
	return "amqp://guest:guest@127.0.0.1:5672/"
}

func Test_putdata(t *testing.T) {
	c.Configuration = c.Configurations{RabbitMQURL: testURL()}
	d, err := New()
	if err != nil {
		t.Fatalf("creating rabbitMQ driver: %v", err)
	}
	if err := d.OpenQueue("teste"); err != nil {
		t.Fatalf("declaring queue teste: %v", err)
	}
	err = d.PutDefault("teste", map[string]string{"msg": "bananinha amassada"})
	if err != nil {
		t.Fatalf("puting data on tube teste: %v", err)
	}
}

func Test_readdata(t *testing.T) {
	c.Configuration = c.Configurations{RabbitMQURL: testURL()}
	d, err := New()
	if err != nil {
		t.Errorf("creating rabbitMQ driver: %v", err)
		return
	}

	if err := d.OpenQueue("teste"); err != nil {
		t.Fatalf("declaring queue teste: %v", err)
	}
	c, err := d.Consume("teste")
	if err != nil {
		t.Fatalf("consuming message from tube teste: %v", err)
	}
	body := <-c
	t.Logf("%s, %s", body.AppId, string(body.Body))
}
