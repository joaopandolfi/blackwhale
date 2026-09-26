package cache

import (
	"time"

	"github.com/joaopandolfi/blackwhale/v2/configurations"
)

const MAX_BUFF_SIZE = 150

var cacheInstance Cache

var InitializedChan chan struct{} = make(chan struct{}, 2)

var waitListenners []chan struct{}

type Cache interface {
	Put(key string, data any, duration time.Duration) error
	Get(key string) (any, error)
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

func AddInitializedListenner(l chan struct{}) {
	if waitListenners == nil {
		waitListenners = []chan struct{}{}
	}
	waitListenners = append(waitListenners, l)
}

func initialized() {
	InitializedChan <- struct{}{}
	for _, c := range waitListenners {
		c <- struct{}{}
	}
	waitListenners = nil
}

func Get() Cache {
	if cacheInstance == nil {
		panic("cache not initialized")
	}
	return cacheInstance
}
