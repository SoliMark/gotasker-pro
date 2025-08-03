package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/service/mock_service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockTaskService(ctrl)
	h := handler.NewTaskHandler(mockSvc)

	router := gin.Default()
	router.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.CreateTask(c)
	})

	t.Run("success", func(t *testing.T) {
		body := `{"title": "New Task", "content": "Detail"}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing title", func(t *testing.T) {
		body := `{"content": "No title"}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockTaskService(ctrl)
	h := handler.NewTaskHandler(mockSvc)

	router := gin.Default()
	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.GetTask(c)
	})

	t.Run("found", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(123)).Return(&model.Task{
			ID:     123,
			UserID: 1,
			Title:  "Task A",
			Status: model.TaskStatusPending,
		}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(999)).Return(nil, nil)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("wrong owner", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(456)).Return(&model.Task{
			ID:     456,
			UserID: 99, // different user
			Title:  "Oops",
		}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/456", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestListTaskss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockTaskService(ctrl)
	h := handler.NewTaskHandler(mockSvc)

	router := gin.Default()
	router.GET("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.ListTasks(c)
	})

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().ListTasks(gomock.Any(), uint(1)).Return([]*model.Task{
			{ID: 1, Title: "T1", Status: model.TaskStatusPending},
			{ID: 2, Title: "T2", Status: model.TaskStatusDone},
		}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
