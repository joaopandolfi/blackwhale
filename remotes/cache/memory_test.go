package cache

import (
	"testing"
	"time"
)

func Test_cache(t *testing.T) {
	Initialize(time.Minute * 1)
	GetMemory().Put("teste", 1234, time.Minute*2)
	val, err := GetMemory().Get("teste")
	if err != nil {
		t.Errorf("get memory error: %v", err)
		return
	}
	t.Log(val)
}
