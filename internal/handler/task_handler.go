package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

type CreateTaskRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type TaskResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type UpdateTaskRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *string `json:"status"` // must be "pending" or "done" if provided
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	task := &model.Task{
		UserID:  userID.(uint),
		Title:   req.Title,
		Content: req.Content,
		Status:  model.TaskStatusPending,
	}

	if err := h.taskService.CreateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, TaskResponse{
		ID:      task.ID,
		Title:   task.Title,
		Content: task.Content,
		Status:  task.Status,
	})
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var taskID uint
	if err := util.ParseUintParam(c, "id", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
	}

	task, err := h.taskService.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	c.JSON(http.StatusOK, TaskResponse{
		ID:      task.ID,
		Title:   task.Title,
		Content: task.Content,
		Status:  task.Status,
	})
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tasks, err := h.taskService.ListTasks(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}
	var res []TaskResponse
	for _, t := range tasks {
		res = append(res, TaskResponse{
			ID:      t.ID,
			Title:   t.Title,
			Content: t.Content,
			Status:  t.Status,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var taskID uint
	if err := util.ParseUintParam(c, "id", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Status != nil &&
		*req.Status != model.TaskStatusPending &&
		*req.Status != model.TaskStatusDone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Content != nil {
		task.Content = *req.Content
	}
	if req.Status != nil {
		task.Status = *req.Status
	}

	if err := h.taskService.UpdateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, TaskResponse{
		ID:      task.ID,
		Title:   task.Title,
		Content: task.Content,
		Status:  task.Status,
	})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid := userIDVal.(uint)

	var taskID uint
	if err := util.ParseUintParam(c, "id", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), uid, taskID); err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, service.ErrPermissionDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		}
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
