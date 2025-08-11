package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/SoliMark/gotasker-pro/internal/service/mock_service"
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
			UserID: 99,
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

func TestUpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockTaskService(ctrl)
	h := handler.NewTaskHandler(mockSvc)

	router := gin.Default()
	router.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.UpdateTask(c)
	})

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(10)).Return(&model.Task{
			ID:     10,
			UserID: 1,
			Title:  "Old",
			Status: model.TaskStatusPending,
		}, nil)

		mockSvc.EXPECT().UpdateTask(gomock.Any(), gomock.AssignableToTypeOf(&model.Task{})).
			Return(nil)

		body := `{"title":"New Title","status":"done"}`
		req, _ := http.NewRequest(http.MethodPut, "/tasks/10", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		body := `{"status":"weird"}`
		req, _ := http.NewRequest(http.MethodPut, "/tasks/10", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(999)).Return(nil, nil)

		req, _ := http.NewRequest(http.MethodPut, "/tasks/999", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc.EXPECT().GetTask(gomock.Any(), uint(456)).Return(&model.Task{
			ID:     456,
			UserID: 777,
			Title:  "X",
		}, nil)

		req, _ := http.NewRequest(http.MethodPut, "/tasks/456", strings.NewReader(`{"title":"abc"}`))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestDeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid id -> 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Set("userID", uint(1))

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized -> 401", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("not found -> 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("userID", uint(1))

		mockSvc.EXPECT().
			DeleteTask(gomock.Any(), uint(1), uint(1)).
			Return(service.ErrTaskNotFound)

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("permission denied -> 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/2", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Set("userID", uint(1))

		mockSvc.EXPECT().
			DeleteTask(gomock.Any(), uint(1), uint(2)).
			Return(service.ErrPermissionDenied)

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("repo/internal error -> 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/3", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "3"}}
		c.Set("userID", uint(1))

		mockSvc.EXPECT().
			DeleteTask(gomock.Any(), uint(1), uint(3)).
			Return(errors.New("DB error"))

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success -> 204 no content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := mock_service.NewMockTaskService(ctrl)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/tasks/10", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "10"}}
		c.Set("userID", uint(1))

		mockSvc.EXPECT().
			DeleteTask(gomock.Any(), uint(1), uint(10)).
			Return(nil)

		h := handler.NewTaskHandler(mockSvc)
		h.DeleteTask(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "", w.Body.String())
	})
}
