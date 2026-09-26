package cache

import (
	"time"

	"github.com/joaopandolfi/blackwhale/v2/utils"
)

type cacheInjectable interface {
	inject(c Cache)
}

// lateInitCache is ised to inject cache on struct after a signal
func lateInitCache(c cacheInjectable) {
	if err := recover(); err == nil {
		return
	}
	utils.Debug("[CACHE][Async loading] waiting for ready cache signal")
	wait := make(chan struct{}, 1)
	AddInitializedListenner(wait)
	go func() {
		ticker := time.NewTicker(time.Second * 40)
		defer ticker.Stop()
		for {
			select {
			case <-wait:
				c.inject(Get())
				utils.Debug("[CACHE][Async loading] cache initialized")
				return
			case <-ticker.C:
				utils.CriticalError("[CACHE][Async loading] still waiting for cache initialization")
			}
		}
	}()
}
