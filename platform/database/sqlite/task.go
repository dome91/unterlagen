package sqlite

import (
	"encoding/json"
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
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.Exec(query,
		task.ID,
		string(task.Type),
		string(task.Status),
		string(task.Payload),
		task.Error,
		task.Attempts,
		task.MaxAttempts,
		task.NextRunAt,
		task.CreatedAt,
		task.UpdatedAt,
	)
	
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
	var taskType, status string
	var payload string
	
	err := r.QueryRow(query, id).Scan(
		&task.ID,
		&taskType,
		&status,
		&payload,
		&task.Error,
		&task.Attempts,
		&task.MaxAttempts,
		&task.NextRunAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err != nil {
		return common.Task{}, err
	}
	
	task.Type = common.TaskType(taskType)
	task.Status = common.TaskStatus(status)
	task.Payload = json.RawMessage(payload)
	
	return task, nil
}

func (r *TaskRepository) FindPendingTasks(limit int) ([]common.Task, error) {
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts, 
			   next_run_at, created_at, updated_at 
		FROM tasks 
		WHERE status = ? AND next_run_at <= ?
		ORDER BY created_at ASC
		LIMIT ?
	`
	
	rows, err := r.Query(query, string(common.TaskStatusPending), time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []common.Task
	for rows.Next() {
		var task common.Task
		var taskType, status string
		var payload string
		
		err := rows.Scan(
			&task.ID,
			&taskType,
			&status,
			&payload,
			&task.Error,
			&task.Attempts,
			&task.MaxAttempts,
			&task.NextRunAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		task.Type = common.TaskType(taskType)
		task.Status = common.TaskStatus(status)
		task.Payload = json.RawMessage(payload)
		
		tasks = append(tasks, task)
	}
	
	return tasks, rows.Err()
}

func (r *TaskRepository) FindAll() ([]common.Task, error) {
	query := `
		SELECT id, type, status, payload, error, attempts, max_attempts, 
			   next_run_at, created_at, updated_at 
		FROM tasks 
		ORDER BY created_at DESC
	`
	
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []common.Task
	for rows.Next() {
		var task common.Task
		var taskType, status string
		var payload string
		
		err := rows.Scan(
			&task.ID,
			&taskType,
			&status,
			&payload,
			&task.Error,
			&task.Attempts,
			&task.MaxAttempts,
			&task.NextRunAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		task.Type = common.TaskType(taskType)
		task.Status = common.TaskStatus(status)
		task.Payload = json.RawMessage(payload)
		
		tasks = append(tasks, task)
	}
	
	return tasks, rows.Err()
}

func (r *TaskRepository) FindPaginated(limit, offset int) ([]common.Task, int, error) {
	// Count total tasks
	var total int
	countQuery := `SELECT COUNT(*) FROM tasks`
	err := r.QueryRow(countQuery).Scan(&total)
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
	
	rows, err := r.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var tasks []common.Task
	for rows.Next() {
		var task common.Task
		var taskType, status string
		var payload string
		
		err := rows.Scan(
			&task.ID,
			&taskType,
			&status,
			&payload,
			&task.Error,
			&task.Attempts,
			&task.MaxAttempts,
			&task.NextRunAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		
		if err != nil {
			return nil, 0, err
		}
		
		task.Type = common.TaskType(taskType)
		task.Status = common.TaskStatus(status)
		task.Payload = json.RawMessage(payload)
		
		tasks = append(tasks, task)
	}
	
	return tasks, total, rows.Err()
}

func (r *TaskRepository) DeleteByID(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.Exec(query, id)
	return err
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}