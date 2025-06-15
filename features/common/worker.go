package common

import (
	"context"
	"log/slog"
	"time"
)

type TaskProcessor interface {
	ProcessTask(task Task) error
	ResponsibleFor() []TaskType
}

type Worker struct {
	repository TaskRepository
	processor  TaskProcessor
	interval   time.Duration
	batchSize  int
}

func NewWorker(repository TaskRepository, processor TaskProcessor) *Worker {
	return &Worker{
		repository: repository,
		processor:  processor,
		interval:   5 * time.Second,
		batchSize:  1,
	}
}

func (w *Worker) Start(ctx context.Context) {
	slog.Info("starting task worker", "interval", w.interval, "batch_size", w.batchSize)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			slog.Info("task worker shutting down")
			return
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
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

func (w *Worker) processTask(task Task) {
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
