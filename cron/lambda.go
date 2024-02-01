package cron

import "context"

type LambdaJob struct {
	Lambda func(ctx context.Context) error
	Ctx    context.Context
}

func (j *LambdaJob) Trigger() bool {
	return true
}

func (j *LambdaJob) Eval() error {
	return j.Lambda(j.Ctx)
}

func (j *LambdaJob) Stop() error {
	return nil
}
