package cache

import (
	"sync"
	"time"

	"github.com/joaopandolfi/blackwhale/v2/configurations"
)

const MAX_BUFF_SIZE = 150

var cacheInstance Cache

var InitializedChan chan struct{} = make(chan struct{}, 2)

var waitListenners []chan struct{}
var waitListennersMu sync.RWMutex

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
	waitListennersMu.Lock()
	waitListenners = append(waitListenners, l)
	waitListennersMu.Unlock()
}

func initialized() {
	waitListennersMu.Lock()
	listeners := waitListenners
	waitListenners = nil
	waitListennersMu.Unlock()

	InitializedChan <- struct{}{}
	for _, c := range listeners {
		c <- struct{}{}
	}
}

func Get() Cache {
	if cacheInstance == nil {
		panic("cache not initialized")
	}
	return cacheInstance
}
