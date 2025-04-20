package service

import (
	"test_task/internal/repository"
	"test_task/internal/service/task"
	"test_task/models"
)

type TaskService interface {
	CreateTask(task *models.Task) (taskId string, err error)
	CheckTaskStatus(taskIdentification string) (*models.Task, error)
	CancelAllTasks()
}

type Service struct {
	TaskService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		TaskService: task.NewTaskService(repos),
	}
}
