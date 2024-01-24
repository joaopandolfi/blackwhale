//go:build integration

package cache

import (
	"testing"
	"time"
)

func Test_redis(t *testing.T) {
	GetRedis().Put("teste", 1234, time.Minute*2)
	val, err := GetRedis().Get("teste")
	if err != nil {
		t.Errorf("get redis error: %v", err)
		return
	}
	t.Log(val)
}
