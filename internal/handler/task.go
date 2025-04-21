package handler

import (
	"encoding/json"
	"net/http"
	"test_task/models"
)

func (h *Handler) tCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	/*
		TODO
		в передаче данных для задачи могли бы быть какие то данные для вычислений,
		но я замокал
	*/

	var input models.Task
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, http.StatusBadRequest, "fail to parse body")
		return
	}

	if input.Task == "" {
		errorResponse(w, http.StatusBadRequest, "task is required")
		return
	}

	taskID, err := h.services.CreateTask(&input)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"Status": "success",
		"TaskId": taskID,
	})
}

func (h *Handler) tStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}

	taskIdentification := r.URL.Query().Get("task_id")
	if taskIdentification == "" {
		errorResponse(w, http.StatusBadRequest, "task_id is required")
		return
	}

	task, err := h.services.CheckTaskStatus(taskIdentification)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Task_identification": task.TaskIdentification,
		"Task":                task.Task,
		"Result":              task.Result,
		"Status":              task.Status,
		"CreatedAt":           task.CreatedAt,
	})
}
