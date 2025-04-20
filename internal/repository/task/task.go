package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"test_task/models"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateTask(request *models.Task) error {
	query := `INSERT INTO task (task_identification, task, result, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, request.TaskIdentification, request.Task, "some result...", "in progress")
	if err != nil {
		return errors.New("Failed to create task")
	}
	return nil
}

func (r *TaskPostgres) TaskResult(taskIdentification string) (*models.Task, error) {
	var output models.Task
	query := `SELECT task_identification, task, result, status, created_at FROM task WHERE task_identification=$1`
	err := r.db.QueryRow(query, taskIdentification).Scan(&output.TaskIdentification, &output.Task, &output.Result, &output.Status, &output.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Task not found")
		}
		return nil, err
	}
	return &output, nil
}

func (r *TaskPostgres) UpdateTaskStatus(taskIdentification string) error {
	query := `UPDATE task SET status=$1 WHERE task_identification=$2`
	_, err := r.db.Exec(query, "completed", taskIdentification)
	if err != nil {
		return errors.New("Failed to update task status")
	}
	return nil
}

func (r *TaskPostgres) TaskFailed(taskIdentification string) error {
	query := `UPDATE task SET status=$1 WHERE task_identification=$2`
	_, err := r.db.Exec(query, "failed", taskIdentification)
	if err != nil {
		return errors.New("Failed to update task status")
	}
	return nil
}

func (r *TaskPostgres) UpdateTasksOnShutdown(ctx context.Context) error {
	query := `
        UPDATE task 
        SET status = 'failed' 
        WHERE status = 'in progress'
    `
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update tasks on shutdown: %w", err)
	}
	return nil
}
