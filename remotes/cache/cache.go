package cache

import (
	"time"

	"github.com/joaopandolfi/blackwhale/configurations"
)

const MAX_BUFF_SIZE = 150

var cacheInstance Cache

var InitializedChan chan bool = make(chan bool, 2)

type Cache interface {
	Put(key string, data interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Size() int
	Flush() error
	GracefullShutdown()
}

func Initialize(tick time.Duration) Cache {
	if configurations.Configuration.Redis.Use {
		cacheInstance = GetRedis()
	} else {
		cacheInstance = initializeMemory(tick)
	}
	InitializedChan <- true
	return cacheInstance
}

func Get() Cache {
	if cacheInstance == nil {
		panic("cache not initialized")
	}
	return cacheInstance
}
