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
