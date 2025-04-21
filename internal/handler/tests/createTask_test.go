package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"test_task/internal/handler"
	"test_task/internal/service"
	serviceMocks "test_task/internal/service/mocks"
	"test_task/models"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type taskServiceMockFunc func(mc *minimock.Controller) service.TaskService

	var (
		mc = minimock.NewController(t)

		task                = gofakeit.StreetName()
		task_identification = gofakeit.UUID()

		req = &models.Task{
			Task: task,
		}

		res = map[string]string{
			"Status": "success",
			"TaskId": task_identification,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		requestBody     *models.Task
		expectedCode    int
		expectedBody    map[string]string
		err             error
		taskServiceMock taskServiceMockFunc
	}{
		{
			name:         "success case",
			requestBody:  req,
			expectedCode: http.StatusCreated,
			expectedBody: res,
			err:          nil,
			taskServiceMock: func(mc *minimock.Controller) service.TaskService {
				mock := serviceMocks.NewTaskServiceMock(mc)
				mock.CreateTaskMock.Expect(req).Return(task_identification, nil)
				return mock
			},
		},
		{
			name:         "task service error",
			requestBody:  &models.Task{Task: ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{"error": "task is required"},
			err:          errors.New("task service error"),
			taskServiceMock: func(mc *minimock.Controller) service.TaskService {
				mock := serviceMocks.NewTaskServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			taskServiceMock := tt.taskServiceMock(mc)

			service := &service.Service{
				TaskService: taskServiceMock,
			}

			api := handler.NewHandler(service)

			jsonData, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/task/create", bytes.NewBuffer(jsonData))
			request.Header.Set("Content-Type", "application/json")

			router := api.RegisterRoutes()
			router.ServeHTTP(recorder, request)

			require.Equal(t, tt.expectedCode, recorder.Code)

			var responseBody map[string]string
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
