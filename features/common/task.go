package common

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
)

type TaskSchedulerMode bool

var TaskSchedulerModeSynchronous TaskSchedulerMode = true
var TaskSchedulerModeAsynchronous TaskSchedulerMode = false

type TaskScheduler struct {
	shutdown       *Shutdown
	taskRepository TaskRepository
	mode           TaskSchedulerMode
}

func (s *TaskScheduler) Register(processor TaskProcessor) {
	worker := newWorker(s.taskRepository, processor)
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddCallback(cancel)
	go worker.start(ctx)
}

func (s *TaskScheduler) ScheduleTask(taskType TaskType, payload any, maxAttempts int) error {
	task, err := newTask(taskType, payload, maxAttempts)
	if err != nil {
		return err
	}

	err = s.taskRepository.Save(task)
	if err != nil {
		return err
	}
	// Create a channel to signal task completion
	resultCh := make(chan error, 1)

	if s.mode == TaskSchedulerModeAsynchronous {
		// Start a goroutine to monitor the task's status
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for range ticker.C {
				currentTask, err := s.taskRepository.FindByID(task.ID)
				if err != nil {
					resultCh <- err
					return
				}

				if currentTask.Status == TaskStatusCompleted {
					resultCh <- nil
					return
				}

				if currentTask.Status == TaskStatusFailed && currentTask.Attempts >= currentTask.MaxAttempts {
					resultCh <- errors.New(currentTask.Error)
					return
				}
			}
		}()

		return <-resultCh
	}

	return nil
}

func NewTaskScheduler(
	shutdown *Shutdown,
	taskRepository TaskRepository,
	mode TaskSchedulerMode,
) *TaskScheduler {
	return &TaskScheduler{
		shutdown:       shutdown,
		taskRepository: taskRepository,
		mode:           mode,
	}
}

type TaskType string
type TaskStatus string

const (
	TaskTypeExtractText      TaskType = "extract_text"
	TaskTypeGeneratePreviews TaskType = "generate_previews"
	TaskTypeIndexDocument    TaskType = "index_document"
)

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string          `json:"id" db:"id"`
	Type        TaskType        `json:"type" db:"type"`
	Status      TaskStatus      `json:"status" db:"status"`
	Payload     json.RawMessage `json:"payload" db:"payload"`
	Error       string          `json:"error,omitempty" db:"error"`
	Attempts    int             `json:"attempts" db:"attempts"`
	MaxAttempts int             `json:"max_attempts" db:"max_attempts"`
	NextRunAt   time.Time       `json:"next_run_at" db:"next_run_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

func newTask(taskType TaskType, payload any, maxAttempts int) (Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:          GenerateID(),
		Type:        taskType,
		Status:      TaskStatusPending,
		Payload:     payloadBytes,
		Attempts:    0,
		MaxAttempts: maxAttempts,
		NextRunAt:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (t *Task) MarkRunning() {
	t.Status = TaskStatusRunning
	t.Attempts++
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkFailed(err error) {
	t.Status = TaskStatusFailed
	t.Error = err.Error()
	t.UpdatedAt = time.Now()

	if t.Attempts < t.MaxAttempts {
		t.Status = TaskStatusPending
		t.scheduleRetry()
	}
}

func (t *Task) scheduleRetry() {
	// TODO: Rather exponential backoff?
	backoffDuration := min(time.Duration(t.Attempts*t.Attempts)*time.Second, 5*time.Minute)
	t.NextRunAt = time.Now().Add(backoffDuration)
}

func (t *Task) ShouldRun() bool {
	return t.Status == TaskStatusPending && time.Now().After(t.NextRunAt)
}

type TaskRepository interface {
	Save(task Task) error
	FindByID(id string) (Task, error)
	FindPendingTasksOfAnyType(limit int, types []TaskType) ([]Task, error)
	FindAll() ([]Task, error)
	FindPaginated(limit, offset int) ([]Task, int, error)
	DeleteByID(id string) error
	DeleteCompleted() error
}

type TaskProcessor interface {
	Name() string
	ProcessTask(task Task) error
	ResponsibleFor() []TaskType
}

type worker struct {
	repository TaskRepository
	processor  TaskProcessor
	interval   time.Duration
	batchSize  int
}

func newWorker(repository TaskRepository, processor TaskProcessor) *worker {
	return &worker{
		repository: repository,
		processor:  processor,
		interval:   5 * time.Millisecond,
		batchSize:  10,
	}
}

func (w *worker) start(ctx context.Context) {
	slog.Info("starting task processor", "name", w.processor.Name(), "interval", w.interval, "batch_size", w.batchSize)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			slog.Info("task processor shutting down")
			return
		}
	}
}

func (w *worker) processBatch(ctx context.Context) {
	tasks, err := w.repository.FindPendingTasksOfAnyType(w.batchSize, w.processor.ResponsibleFor())
	if err != nil {
		slog.Error("failed to find pending tasks", "error", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	slog.Debug("processing task batch", "count", len(tasks))

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			w.processTask(task)
		}
	}
}

func (w *worker) processTask(task Task) {
	if !task.ShouldRun() {
		return
	}

	slog.Info("processing task", "id", task.ID, "type", task.Type, "attempt", task.Attempts+1)

	task.MarkRunning()
	err := w.repository.Save(task)
	if err != nil {
		slog.Error("failed to mark task as running", "task_id", task.ID, "error", err)
		return
	}

	err = w.processor.ProcessTask(task)
	if err != nil {
		slog.Error("task failed", "task_id", task.ID, "error", err, "attempt", task.Attempts)
		task.MarkFailed(err)
	} else {
		slog.Info("task completed", "task_id", task.ID, "type", task.Type)
		task.MarkCompleted()
	}

	err = w.repository.Save(task)
	if err != nil {
		slog.Error("failed to save task status", "task_id", task.ID, "error", err)
	}
}
