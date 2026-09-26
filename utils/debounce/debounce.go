package debounce

import (
	"sync"
	"time"

	"github.com/joaopandolfi/blackwhale/v2/utils"
)

var (
	debounces    = sync.Map{}
	channels     = sync.Map{}
	counter      = map[string]int{}
	counterMutex sync.Mutex
)

const (
	channelBuffer = 200
)

// Create and returns a channel bound to given id
func Channel(id string) chan any {
	ch, _ := channels.LoadOrStore(id, make(chan any, channelBuffer))
	return ch.(chan any)
}

// Runs the callback passing input as payload when interval is over. Reset interval whenever input channel receives a new payload.
func Run(
	id string,
	interval time.Duration,
	input chan any,
	callback func(payload any),
	logMsg_optional ...string,
) {
	// skip execution if debounce is already running for given id
	if active, ok := debounces.Load(id); ok && active.(bool) {
		return
	}

	setupDebounce(id)

	silent := len(logMsg_optional) == 0

	var payload any
	timer := time.NewTimer(interval)
	for {
		select {
		case payload = <-input:
			timer.Reset(interval)
			increaseCounter(id)
		case <-timer.C:
			go callback(payload)
			if !silent {
				log(id, logMsg_optional[0])
			}
			clear(id)
			return
		}
	}
}

func log(id, logMsg_optional string) {
	utils.Debug("[Debounce] - Buffered calls", counter[id], logMsg_optional)
}

func setupDebounce(id string) {
	debounces.Store(id, true)
	setCounter(id)
}

func clearDebounce(id string) {
	debounces.Delete(id)
}

func clear(id string) {
	clearDebounce(id)
	clearChannel(id)
	clearCounter(id)
}

func clearChannel(id string) {
	channels.Delete(id)
}

func setCounter(id string) {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	counter[id] = 0
}

func clearCounter(id string) {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	delete(counter, id)

}

func increaseCounter(id string) int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	counter[id]++
	return counter[id]
}
