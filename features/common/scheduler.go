package common

import "context"

type JobScheduler struct {
	shutdown       *Shutdown
	taskRepository TaskRepository
}

func (s *JobScheduler) Schedule(task func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go task(ctx)
}

func NewJobScheduler(shutdown *Shutdown) *JobScheduler {
	return &JobScheduler{
		shutdown: shutdown,
	}
}

type TaskScheduler struct {
	shutdown       *Shutdown
	taskRepository TaskRepository
}

func (s *TaskScheduler) RegisterWorker(processor TaskProcessor) {
	worker := NewWorker(s.taskRepository, processor)
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go worker.Start(ctx)
}

func (s *TaskScheduler) ScheduleTask(taskType TaskType, payload any, maxAttempts int) error {
	task, err := NewTask(taskType, payload, maxAttempts)
	if err != nil {
		return err
	}

	return s.taskRepository.Save(task)
}

func NewTaskScheduler(
	shutdown *Shutdown,
	taskRepository TaskRepository) *TaskScheduler {
	return &TaskScheduler{
		shutdown:       shutdown,
		taskRepository: taskRepository,
	}
}
