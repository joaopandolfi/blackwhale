package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaopandolfi/blackwhale/utils"
)

type work struct {
	job    Job
	tick   time.Duration
	active bool
	stop   chan bool
}

type cron struct {
	jobs          map[string]*work
	mu            sync.Mutex
	onlineWorkers int
	stopCh        chan bool
	endCh         chan bool
	errCh         chan error
}

var cr *cron

func Get() CRON {
	if cr == nil {
		cr = &cron{
			jobs: map[string]*work{},
		}
	}
	return cr
}

func (c *cron) AddJob(key string, tick time.Duration, job Job) error {
	c.init()

	if c.jobs[key] != nil {
		return fmt.Errorf("job (%s) already registered", key)
	}

	c.jobs[key] = &work{
		job:    job,
		tick:   tick,
		active: true,
		stop:   make(chan bool),
	}

	return nil
}

func (c *cron) RemoveJob(key string) error {
	c.init()

	if c.jobs[key] == nil {
		return fmt.Errorf("job (%s) does not exits", key)
	}

	if c.jobs[key].active {
		return fmt.Errorf("a running job cant be deleted")
	}

	delete(c.jobs, key)
	return nil
}

func (c *cron) StopJob(key string) error {
	c.init()

	if c.jobs[key] == nil {
		return fmt.Errorf("job (%s) does not exits", key)
	}

	c.jobs[key].stop <- true

	return nil
}

func (c *cron) Start() {
	c.stopCh = make(chan bool)
	c.endCh = make(chan bool)
	c.errCh = make(chan error, len(c.jobs))

	for k, j := range c.jobs {
		go c.worker(k, j.stop, j.tick, j.job)
	}

	go c.errorHandler()
}

func (c *cron) GracefullShutdown() error {
	c.stopCh <- true
	<-c.endCh
	return nil
}

func (c *cron) errorHandler() {
	c.workerStarted()
	for {
		select {
		case <-c.stopCh:
			utils.Info("[CRON][Stop] handler", "error")
			c.stopPropagate()
			return
		case err := <-c.errCh:
			fmt.Println("[CROTN][Error Handler] Error: %w", err)
		}
	}
}

func (c *cron) worker(key string, stop chan bool, tick time.Duration, job Job) {
	c.workerStarted()
	ticker := time.NewTicker(tick)

	utils.Info("[CRON][Start] job", key, tick/time.Second, "seconds")

	for {
		select {
		case <-c.stopCh:
			utils.Info("[CRON][Stop] job", key)
			c.stopPropagate()
			return
		case <-stop:
			utils.Info("[CRON][Stop] only job", key)
			c.workerStopped(key)
			return
		case <-ticker.C:
			if job.Trigger() {
				err := job.Eval()
				if err != nil {
					c.errCh <- fmt.Errorf("error on job eval: %w", err)
				}
			}
		}
	}
}

func (c *cron) init() {
	if c.jobs == nil {
		c.jobs = map[string]*work{}
	}
}

func (c *cron) workerStarted() {
	c.mu.Lock()
	c.onlineWorkers += 1
	c.mu.Unlock()
}

func (c *cron) workerStopped(key string) {
	c.mu.Lock()
	c.onlineWorkers -= 1
	c.mu.Unlock()
	if c.jobs[key] != nil {
		c.jobs[key].active = false
	}
}

func (c *cron) stopPropagate() {
	c.onlineWorkers--
	if c.onlineWorkers <= 0 {
		c.endCh <- true
		return
	}
	time.Sleep(5 * time.Millisecond)
	c.stopCh <- true
}
