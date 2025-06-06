package common

import "context"

type Scheduler struct {
	shutdown *Shutdown
}

func (s *Scheduler) Schedule(task func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go task(ctx)
}

func NewScheduler(shutdown *Shutdown) *Scheduler {
	return &Scheduler{shutdown: shutdown}
}
