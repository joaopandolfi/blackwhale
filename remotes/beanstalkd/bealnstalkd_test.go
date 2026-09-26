//go:build integration

package beanstalkd

import (
	"os"
	"testing"

	c "github.com/joaopandolfi/blackwhale/v2/configurations"
)

func testURL() string {
	if v := os.Getenv("BEANSTALKD_URL"); v != "" {
		return v
	}
	return "127.0.0.1:11300"
}

func Test_putdata(t *testing.T) {
	c.Configuration = c.Configurations{BeanstalkdUrl: testURL()}
	b, err := New()
	if err != nil {
		t.Errorf("Creating beanstalkd: %v", err)
		return
	}
	_, err = b.PutDefault("/test/blackwale", map[string]string{"message": "Testing1"})
	if err != nil {
		t.Errorf("Putting data on tube: %v", err)
	}
}

func Test_readdata(t *testing.T) {
	c.Configuration = c.Configurations{BeanstalkdUrl: testURL()}
	b, err := New()
	if err != nil {
		t.Errorf("Creating beanstalkd: %v", err)
		return
	}
	id, bytes, err := b.ReadTube("/test/blackwale", DURATION_DEFAULT)
	if err != nil {
		t.Errorf("Reading data from tube: %v", err)
		return
	}

	t.Logf("received data: %s", string(bytes))
	err = b.DeleteMessage("/test/blackwale", id)
	if err != nil {
		t.Errorf("Deleting message: %v", err)
	}
}
