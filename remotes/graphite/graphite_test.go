//go:build integration

package graphite

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_general(t *testing.T) {
	host := "127.0.0.1"
	if v := os.Getenv("GRAPHITE_HOST"); v != "" {
		host = v
	}
	port := 2003
	if v := os.Getenv("GRAPHITE_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	SetCredentials(host, port)
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
