package cache

import (
	"time"

	"github.com/joaopandolfi/blackwhale/configurations"
)

const MAX_BUFF_SIZE = 150

var cacheInstance Cache

var InitializedChan chan bool = make(chan bool, 2)

var waitListenners []chan bool

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
	initialized()
	return cacheInstance
}

func AddInitializedListenner(l chan bool) {
	if waitListenners == nil {
		waitListenners = []chan bool{}
	}
	waitListenners = append(waitListenners, l)
}

func initialized() {
	InitializedChan <- true
	for _, c := range waitListenners {
		c <- true
	}
	waitListenners = nil
}

func Get() Cache {
	if cacheInstance == nil {
		panic("cache not initialized")
	}
	return cacheInstance
}
