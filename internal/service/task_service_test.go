package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository/mock_repository"
	"github.com/SoliMark/gotasker-pro/internal/service"
)

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)
	svc := service.NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		task := &model.Task{UserID: 1, Title: "My Task"}
		mockRepo.EXPECT().CreateTask(ctx, task).Return(nil)

		err := svc.CreateTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("empty title", func(t *testing.T) {
		task := &model.Task{UserID: 2, Title: ""}
		err := svc.CreateTask(ctx, task)
		assert.EqualError(t, err, "title is required")
	})

	t.Run("repo returns error", func(t *testing.T) {
		task := &model.Task{UserID: 3, Title: "Fail Task"}
		mockRepo.EXPECT().CreateTask(ctx, task).Return(errors.New("DB error"))

		err := svc.CreateTask(ctx, task)
		assert.EqualError(t, err, "DB error")
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)
	svc := service.NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		task := &model.Task{
			ID:     10,
			UserID: 1,
			Title:  "New Title",
			Status: model.TaskStatusDone,
		}
		mockRepo.EXPECT().UpdateTask(ctx, task).Return(nil)

		err := svc.UpdateTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("empty title -> error", func(t *testing.T) {
		task := &model.Task{
			ID:     10,
			UserID: 1,
			Title:  "",
		}
		err := svc.UpdateTask(ctx, task)
		assert.EqualError(t, err, "title is required")
	})

	t.Run("repo error", func(t *testing.T) {
		task := &model.Task{
			ID:     10,
			UserID: 1,
			Title:  "X",
		}
		mockRepo.EXPECT().UpdateTask(ctx, task).Return(errors.New("db err"))

		err := svc.UpdateTask(ctx, task)
		assert.EqualError(t, err, "db err")
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)
	svc := service.NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, uint(10)).
			Return(&model.Task{ID: 10, UserID: 1}, nil)
		mockRepo.EXPECT().DeleteTask(ctx, uint(10)).Return(nil)

		err := svc.DeleteTask(ctx, 1, 10) // 傳 userID=1, taskID=10
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, uint(11)).
			Return(nil, nil)

		err := svc.DeleteTask(ctx, 1, 11)
		assert.EqualError(t, err, service.ErrTaskNotFound.Error())
	})

	t.Run("permission denied", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, uint(12)).
			Return(&model.Task{ID: 12, UserID: 99}, nil) // 任務屬於別人

		err := svc.DeleteTask(ctx, 1, 12)
		assert.EqualError(t, err, service.ErrPermissionDenied.Error())
	})

	t.Run("repo returns error", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, uint(13)).
			Return(&model.Task{ID: 13, UserID: 1}, nil)
		mockRepo.EXPECT().DeleteTask(ctx, uint(13)).
			Return(errors.New("DB error"))

		err := svc.DeleteTask(ctx, 1, 13)
		assert.EqualError(t, err, "DB error")
	})
}
