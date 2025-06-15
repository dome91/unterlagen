package sqlite

import (
	"time"
	"unterlagen/features/common"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	*sqlx.DB
}

func (r *TaskRepository) Save(task common.Task) error {
	query := `
		INSERT OR REPLACE INTO tasks (
			id, type, status, payload, error, attempts, max_attempts,
			next_run_at, created_at, updated_at
		) VALUES (:id, :type, :status, :payload, :error, :attempts, :max_attempts,
			:next_run_at, :created_at, :updated_at)
	`

	_, err := r.NamedExec(query, task)
	return err
}

func (r *TaskRepository) FindByID(id string) (common.Task, error) {
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts,
			   next_run_at, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	var task common.Task
	err := r.Get(&task, query, id)
	if err != nil {
		return common.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) FindPendingTasksOfAnyType(limit int, types []common.TaskType) ([]common.Task, error) {
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts,
			   next_run_at, created_at, updated_at
		FROM tasks
		WHERE status = ? AND next_run_at <= ? AND type IN (?)
		ORDER BY created_at ASC
		LIMIT ?
	`

	var typesArg []string
	for _, t := range types {
		typesArg = append(typesArg, string(t))
	}
	query, args, err := sqlx.In(query, string(common.TaskStatusPending), time.Now(), typesArg, limit)
	if err != nil {
		return nil, err
	}

	query = r.Rebind(query)
	var tasks []common.Task
	err = r.Select(&tasks, query, args...)
	return tasks, err
}

func (r *TaskRepository) FindAll() ([]common.Task, error) {
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts,
			   next_run_at, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	var tasks []common.Task
	err := r.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) FindPaginated(limit, offset int) ([]common.Task, int, error) {
	// Count total tasks
	var total int
	countQuery := `SELECT COUNT(*) FROM tasks`
	err := r.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated tasks
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts,
			   next_run_at, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	var tasks []common.Task
	err = r.Select(&tasks, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) DeleteByID(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.Exec(query, id)
	return err
}

func (r *TaskRepository) DeleteCompleted() error {
	query := `DELETE FROM tasks WHERE status = ?`
	_, err := r.Exec(query, string(common.TaskStatusCompleted))
	return err
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}
