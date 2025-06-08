package common

import "context"

type Scheduler struct {
	shutdown       *Shutdown
	taskRepository TaskRepository
}

func (s *Scheduler) Schedule(task func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go task(ctx)
}

func (s *Scheduler) RegisterWorker(processor TaskProcessor) {
	worker := NewWorker(s.taskRepository, processor)
	s.Schedule(worker.Start)
}

func (s *Scheduler) ScheduleTask(taskType TaskType, payload any, maxAttempts int) error {
	task, err := NewTask(taskType, payload, maxAttempts)
	if err != nil {
		return err
	}

	return s.taskRepository.Save(task)
}

func (s *Scheduler) SetTaskRepository(taskRepository TaskRepository) {
	s.taskRepository = taskRepository
}

func NewScheduler(shutdown *Shutdown) *Scheduler {
	return &Scheduler{
		shutdown: shutdown,
	}
}
