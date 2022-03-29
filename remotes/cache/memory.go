package cache

import (
	"sync"
	"time"

	"github.com/joaopandolfi/blackwhale/utils"
	"golang.org/x/xerrors"
)

var mcache *memCache

type memCache struct {
	buff        map[string]stored
	garbageStop chan bool
	mu          sync.RWMutex
}

func GetMemory() Cache {
	return mcache
}

func initializeMemory(tick time.Duration) Cache {
	if mcache == nil {
		utils.Info("[CACHE] using local cache", "Memory")
		mcache = &memCache{
			buff: map[string]stored{},
		}

		mcache.startGarbageCollector(tick)
	}
	return mcache
}

func (c *memCache) Put(key string, data interface{}, duration time.Duration) error {
	if len(c.buff) > MAX_BUFF_SIZE {
		return xerrors.Errorf("buffer overflow")
	}

	c.buff[key] = stored{
		value:   data,
		validAt: time.Now().Add(duration),
	}
	return nil
}

func (c *memCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.buff[key]; ok {
		return val.value, nil
	}
	return nil, nil
}

func (c *memCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.buff, key)

	return nil
}

func (c *memCache) Flush() error {
	c.mu.Lock()
	c.buff = map[string]stored{}
	c.mu.Unlock()
	return nil
}

func (c *memCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.buff)
}

// ==== GARBAGE COLLECTOR -> PUT THIS IN A SEPARATED STRUCTURE

func (c *memCache) startGarbageCollector(tick time.Duration) {
	ticker := time.NewTicker(tick)
	c.garbageStop = make(chan bool)

	go func() {
		utils.Info("[LOCAL_CACHE][GARBAGE COLLECTOR]", "START")
		for {
			select {
			case <-c.garbageStop:
				ticker.Stop()
				utils.Info("[LOCAL_CACHE][GARBAGE COLLECTOR]", "STOP")
				return
			case <-ticker.C:
				c.GarbageCollector()
			}
		}
	}()
}

func (c *memCache) GarbageCollector() {
	var toDelete []string
	for k, val := range c.buff {
		if val.validAt.After(time.Now()) {
			toDelete = append(toDelete, k)
		}
	}

	for _, d := range toDelete {
		c.Delete(d)
	}
}

func (c *memCache) GracefullShutdown() {
	if c.garbageStop != nil {
		c.garbageStop <- true
	}
}
