package task

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"test_task/internal/repository"
	"test_task/models"
)

func generateUniqueID() string {
	return uuid.New().String()
}

type TaskService struct {
	repo  repository.TaskRepository
	tasks sync.Map
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *models.Task) (taskId string, err error) {
	task.TaskIdentification = generateUniqueID()

	err = s.repo.CreateTask(task)
	if err != nil {
		return "", err
	}

	taskCtx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	s.tasks.Store(task.TaskIdentification, cancel)

	go func() {
		defer cancel()
		s.runTask(taskCtx, task.TaskIdentification)
	}()
	return task.TaskIdentification, nil
}

func (s *TaskService) CheckTaskStatus(taskIdentification string) (*models.Task, error) {
	task, err := s.repo.TaskResult(taskIdentification)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) CancelAllTasks() {
	s.tasks.Range(func(key, value interface{}) bool {
		cancel := value.(context.CancelFunc)
		cancel()
		return true
	})
}

func (s *TaskService) runTask(ctx context.Context, taskIdentification string) {
	select {
	case <-ctx.Done():
		log.Println("Task execution canceled:", ctx.Err())
		err := s.repo.TaskFailed(taskIdentification)
		if err != nil {
			log.Println("Failed to mark task as failed:", err)
		}
		return
	default:
		time.Sleep(3 * time.Minute) // Имитация какой-нибудь долгосрочной задачи

		/*
			TODO
			здесь должен сохраняться результат работы определенной задачи,
			но на данный момент я замокал данные еще при их отправке в CreateTask
		*/
		err := s.repo.UpdateTaskStatus(taskIdentification)
		if err != nil {
			log.Println("Failed to update task status:", err)
		}
		return
	}
}
