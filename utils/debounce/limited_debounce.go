package debounce

import "time"

// Runs the callback passing input as payload when interval is over. Reset interval whenever input channel receives a new payload.
func RunLimited(
	id string,
	interval time.Duration,
	input chan interface{},
	limit int,
	callback func(payload interface{}),
	logMsg_optional ...string,
) {
	// skip execution if debounce is already running for given id
	if active, ok := debounces.Load(id); ok && active.(bool) {
		return
	}

	setupDebounce(id)

	silent := len(logMsg_optional) == 0

	var payload interface{}
	timer := time.NewTimer(interval)
	tick := make(chan bool)
	for {
		select {
		case payload = <-input:
			timer.Reset(interval)
			if increaseCounter(id) >= limit {
				timer.Stop()
				tick <- true
			}
		case <-timer.C:
			tick <- true
		case <-tick:
			go callback(payload)
			if !silent {
				log(id, logMsg_optional[0])
			}
			clear(id)
			return
		}
	}
}
