package common

import "context"

type Job func(ctx context.Context)

type JobScheduler struct {
	shutdown *Shutdown
}

func (s *JobScheduler) Schedule(job Job) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go job(ctx)
}

func NewJobScheduler(shutdown *Shutdown) *JobScheduler {
	return &JobScheduler{
		shutdown: shutdown,
	}
}
