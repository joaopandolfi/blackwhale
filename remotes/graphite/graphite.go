package graphite

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	conf "github.com/joaopandolfi/blackwhale/configurations"
)

var _host string
var _port int

var _counter map[string]int
var _histogram map[string]*histogram

var _queue chan payload
var _seconds int
var _active bool

var mu sync.Mutex
var muH sync.Mutex

type histogram struct {
	Total int
	Mean  int
}

type payload struct {
	Key string
	Val string
}

// Driver graphite
type Driver struct {
	Conn   *Graphite
	Prefix string
}

var conn *Graphite

func init() {
	_host = conf.Configuration.GraphiteUrl
	p, _ := strconv.Atoi(conf.Configuration.GraphitePort)
	_port = p
	_counter = map[string]int{}
	_histogram = map[string]*histogram{}
	_queue = make(chan payload, 100)
	_seconds = 5
}

// SetCredentials to connection
func SetCredentials(host string, port int) {
	_host = host
	_port = port
}

// SetSeconds to buffer send data
func SetSeconds(seconds int) {
	_seconds = seconds
}

// New Graphite driver
func New(prefix string) (*Driver, error) {
	if conn == nil {
		c, err := NewGraphite(_host, _port)
		if err != nil {
			c = NewGraphiteNop(_host, _port)
		}
		conn = c
	}
	dr := &Driver{
		Conn:   conn,
		Prefix: prefix,
	}

	go dr.sender()
	go dr.flusher(_seconds)

	return dr, nil
}

// Send to graphite
func (d *Driver) Send(key, data string) error {
	return d.Conn.SimpleSend(fmt.Sprintf("stats.%s.%s", d.Prefix, key), data)
}

// Count metric
func (d *Driver) Count(key string) {
	mu.Lock()
	_counter[key]++
	mu.Unlock()
}

// Compute the mean value
func (d *Driver) Compute(key string, value int) {
	muH.Lock()
	h := _histogram[key]
	if h == nil {
		h = &histogram{
			Total: 0,
			Mean:  0,
		}
	}

	h.Total++
	h.Mean = (h.Mean + value) / h.Total

	_histogram[key] = h
	muH.Unlock()
}

// Shutdown - For use on graceful shutdown
func (d *Driver) Shutdown() {
	_active = false
	d.flush()
}

func (d *Driver) sender() {
	for {
		payload := <-_queue
		d.Conn.SimpleSend(fmt.Sprintf("stats.%s.%s", d.Prefix, payload.Key), payload.Val)
	}
}

func (d *Driver) flush() {
	var buff int
	for k := range _counter {
		if k == "" {
			break
		}
		mu.Lock()
		buff = _counter[k]
		_counter[k] = 0
		mu.Unlock()

		if buff == 0 {
			continue
		}
		_queue <- payload{
			Key: k,
			Val: fmt.Sprint(buff),
		}
	}

	for k := range _histogram {
		if k == "" {
			break
		}
		muH.Lock()
		buff = _histogram[k].Mean
		_histogram[k].Total = 0
		_histogram[k].Mean = 0
		muH.Unlock()

		if buff == 0 {
			continue
		}
		_queue <- payload{
			Key: k,
			Val: fmt.Sprint(buff),
		}
	}

}

func (d *Driver) flusher(seconds int) {
	_active = true
	for _active {
		time.Sleep(time.Duration(seconds) * time.Second)
		d.flush()
	}
}
