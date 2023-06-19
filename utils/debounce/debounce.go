package debounce

import (
	"sync"
	"time"

	"github.com/joaopandolfi/blackwhale/utils"
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
func Channel(id string) chan interface{} {
	ch, _ := channels.LoadOrStore(id, make(chan interface{}, channelBuffer))
	return ch.(chan interface{})
}

// Runs the callback passing input as payload when interval is over. Reset interval whenever input channel receives a new payload.
func Run(
	id string,
	interval time.Duration,
	input chan interface{},
	callback func(payload interface{}),
) {
	// skip execution if debounce is already running for given id
	if active, ok := debounces.Load(id); ok && active.(bool) {
		return
	}

	setupDebounce(id)

	var payload interface{}
	timer := time.NewTimer(interval)
	for {
		select {
		case payload = <-input:
			timer.Reset(interval)
			increaseCounter(id)
		case <-timer.C:
			go callback(payload)
			log(id)
			clear(id)
			return
		}
	}
}

func log(id string) {
	utils.Debug("[Debounce] - Buffered calls", counter[id])
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
	counter[id] = 0
	counterMutex.Unlock()
}

func clearCounter(id string) {
	counterMutex.Lock()
	delete(counter, id)
	counterMutex.Unlock()
}

func increaseCounter(id string) {
	counterMutex.Lock()
	counter[id]++
	counterMutex.Unlock()
}
