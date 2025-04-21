package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"test_task/internal/handler"
	"test_task/internal/service"
	serviceMocks "test_task/internal/service/mocks"
	"test_task/models"
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	type taskServiceMockFunc func(mc *minimock.Controller) service.TaskService

	var (
		mc = minimock.NewController(t)

		task                = gofakeit.StreetName()
		task_identification = gofakeit.UUID()
		result              = gofakeit.StreetName()
		status              = gofakeit.Email()
		created_at          = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

		req = task_identification

		res = &models.Task{
			TaskIdentification: task_identification,
			Task:               task,
			Result:             result,
			Status:             status,
			CreatedAt:          created_at,
		}
	)
	defer t.Cleanup(mc.Finish)

	var tests = []struct {
		name            string
		requestBody     string
		expectedCode    int
		expectedBody    *models.Task
		err             error
		taskServiceMock taskServiceMockFunc
	}{
		{
			name:         "success",
			requestBody:  req,
			expectedCode: http.StatusOK,
			expectedBody: res,
			err:          nil,
			taskServiceMock: func(mc *minimock.Controller) service.TaskService {
				mock := serviceMocks.NewTaskServiceMock(mc)
				mock.CheckTaskStatusMock.Expect(req).Return(res, nil)
				return mock
			},
		},
		{
			name:         "failure",
			requestBody:  "",
			expectedCode: http.StatusInternalServerError,
			expectedBody: nil,
			err:          errors.New("internal server error"),
			taskServiceMock: func(mc *minimock.Controller) service.TaskService {
				mock := serviceMocks.NewTaskServiceMock(mc)
				mock.CheckTaskStatusMock.Expect(req).Return(nil, errors.New("internal server error"))
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

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/task/check-result?task_id="+task_identification, nil)
			request.Header.Set("Content-Type", "application/json")

			router := api.RegisterRoutes()
			router.ServeHTTP(recorder, request)

			require.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedBody != nil {
				var responseBody models.Task
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				fmt.Printf("Expected CreatedAt: %v\n", tt.expectedBody.CreatedAt)
				fmt.Printf("Actual CreatedAt: %v\n", responseBody.CreatedAt)
				require.Equal(t, tt.expectedBody, &responseBody)
			}
		})
	}
}
