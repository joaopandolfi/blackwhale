package cron

import (
	"time"
)

type Job interface {
	Trigger() bool
	Eval() error
}

type CRON interface {
	Start()
	AddJob(key string, tick time.Duration, job Job) error
	RemoveJob(key string) error
	StopJob(key string) error
	GracefullShutdown() error
}
