package common

import (
	"encoding/json"
	"time"
)

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

func NewTask(taskType TaskType, payload any, maxAttempts int) (Task, error) {
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
