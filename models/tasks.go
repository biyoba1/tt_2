package models

import "time"

type Task struct {
	ID                 int       `json:"id"`
	TaskIdentification string    `json:"task_identification"` // уникальная uuid для каждой задачи
	Task               string    `json:"task"`                // название задачи
	Result             string    `json:"result"`              // результат выполнения задачи
	Status             string    `json:"status"`              // статус задачи
	CreatedAt          time.Time `json:"created_at"`          // время создания задачи (может понадобится для метрик)
}
