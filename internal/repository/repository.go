package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	task "test_task/internal/repository/task"
	"test_task/models"
)

type TaskRepository interface {
	CreateTask(request *models.Task) error
	TaskResult(taskIdentification string) (*models.Task, error)
	UpdateTaskStatus(taskName string) error
	TaskFailed(taskIdentification string) error
	UpdateTasksOnShutdown(ctx context.Context) error
}

type Repository struct {
	TaskRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TaskRepository: task.NewTaskPostgres(db),
	}
}
