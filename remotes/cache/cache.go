package cache

import "time"

const MAX_BUFF_SIZE = 150

type Cache interface {
	Put(key string, data interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Size() int
	Flush() error
	GracefullShutdown()
}

func Initialize(tick time.Duration) Cache {
	return initializeMemory(tick)
}
