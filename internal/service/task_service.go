package service

import (
	"context"
	"errors"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTask(ctx context.Context, id uint) (*model.Task, error)
	ListTasks(ctx context.Context, userID uint) ([]*model.Task, error)
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, task *model.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}
	return s.repo.CreateTask(ctx, task)
}

func (s *taskService) GetTask(ctx context.Context, id uint) (*model.Task, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *taskService) ListTasks(ctx context.Context, userID uint) ([]*model.Task, error) {
	return s.repo.ListByUserID(ctx, userID)
}
