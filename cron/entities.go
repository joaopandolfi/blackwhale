package cron

import (
	"time"
)

type Job interface {
	Trigger() bool
	Eval() error
	Stop() error
}

type CRON interface {
	Start()
	AddJob(key string, tick time.Duration, job Job) error
	// AddEphemeralJob - Add a job who will run once after a given time
	AddEphemeralJob(key string, tick time.Duration, job Job) error
	RemoveJob(key string) error
	StopJob(key string) error
	GracefullShutdown() error
}

type MockJob struct {
	Ran          bool
	TriggerCount int
	EvalCount    int
	Stopped      bool
}

func (m *MockJob) Trigger() bool {
	m.TriggerCount++
	return !m.Ran
}

func (m *MockJob) Eval() error {
	m.Ran = true
	m.EvalCount++
	return nil
}

func (m *MockJob) Stop() error {
	m.Stopped = true
	return nil
}
