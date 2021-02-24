package graphite

import (
	"fmt"
	"testing"
	"time"
)

func Test_general(t *testing.T) {
	SetCredentials("localhost", 2003)
	d, err := New("test")
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
		return
	}
	go d.Count("banana")
	go d.Count("banana")
	go d.Count("banana")
	go d.Count("banana2")
	go d.Count("banana")

	time.Sleep(10 * time.Second)
}
