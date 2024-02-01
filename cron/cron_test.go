package cron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCron(t *testing.T) {
	cron := &cron{
		jobs: map[string]*work{},
	}
	jobName := "mock"

	mockedJob := MockJob{}
	cron.AddJob(jobName, 1*time.Second, &mockedJob)
	cron.Start()

	time.Sleep(2 * time.Second)
	if !mockedJob.Ran {
		t.Errorf("job not ran")
	}

	err := cron.StopJob(jobName)
	if err != nil {
		t.Errorf("stopping mock job: %v", err)
	}

	time.Sleep(1 * time.Second)
	if !mockedJob.Stopped {
		t.Errorf("job not stopped")
	}

	if cron.onlineWorkers > 1 {
		t.Errorf("Still have online workers: %d", cron.onlineWorkers)
	}

	err = cron.RemoveJob(jobName)
	if err != nil {
		t.Errorf("removing jobs: %v", err)
	}

	if len(cron.jobs) > 0 {
		t.Errorf("Still have registered jobs: %v", cron.jobs)
	}
}

func TestEphemeralJobs(t *testing.T) {
	cron := &cron{
		jobs: map[string]*work{},
	}
	jobName := "mock"

	mockedJob := MockJob{}
	cron.AddEphemeralJob(jobName, EphemeralJobDefaultDelay, &mockedJob)
	cron.Start()
	time.Sleep(2 * time.Second)
	if !mockedJob.Ran {
		t.Errorf("job not ran")
	}
	err := cron.StopJob(jobName)
	if err != nil {
		t.Errorf("stopping mock job: %v", err)
	}

	time.Sleep(1 * time.Second)
	if !mockedJob.Stopped {
		t.Errorf("job not stopped")
	}

	if mockedJob.TriggerCount != 1 {
		t.Errorf("mocked job triggered different from 1 time: %d", mockedJob.TriggerCount)
	}

}

func TestLambda(t *testing.T) {
	cron := &cron{
		jobs: map[string]*work{},
	}

	executed := false

	cron.AddEphemeralJob("lambda", EphemeralJobDefaultDelay, &LambdaJob{
		Ctx: context.Background(),
		Lambda: func(ctx context.Context) error {
			executed = true
			return nil
		},
	})
	cron.Start()

	time.Sleep(time.Second)

	require.Equal(t, executed, true)
}
