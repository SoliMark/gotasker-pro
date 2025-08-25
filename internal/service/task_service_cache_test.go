package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SoliMark/gotasker-pro/internal/cache"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository/mock_repository"
	"github.com/SoliMark/gotasker-pro/internal/service"
)

func TestTaskService_ListTasks_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	// Test data
	userID := uint(1)
	expectedTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: model.TaskStatusPending},
		{ID: 2, UserID: userID, Title: "Task 2", Status: model.TaskStatusDone},
	}

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	// Expect repository call on cache miss
	mockRepo.EXPECT().
		ListByUserID(gomock.Any(), userID).
		Return(expectedTasks, nil)

	tasks, err := service.ListTasks(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)

	// Verify cache was set
	key := cache.KeyUserTasks(userID)
	cachedData, err := rdb.Get(context.Background(), key).Bytes()
	require.NoError(t, err)

	var cachedTasks []*model.Task
	err = json.Unmarshal(cachedData, &cachedTasks)
	require.NoError(t, err)
	assert.Equal(t, expectedTasks, cachedTasks)
}

func TestTaskService_ListTasks_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	// Test data
	userID := uint(1)
	expectedTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: model.TaskStatusPending},
		{ID: 2, UserID: userID, Title: "Task 2", Status: model.TaskStatusDone},
	}

	// Pre-populate cache
	key := cache.KeyUserTasks(userID)
	cachedData, err := json.Marshal(expectedTasks)
	require.NoError(t, err)
	err = rdb.Set(context.Background(), key, cachedData, 60*time.Second).Err()
	require.NoError(t, err)

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	// Should not call repository when cache exists
	// (no mock expectations set)

	tasks, err := service.ListTasks(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)
}

func TestTaskService_ListTasks_WithoutCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// No Redis client (cache disabled)
	service := service.NewTaskService(mockRepo, nil, 60*time.Second)

	userID := uint(1)
	expectedTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: model.TaskStatusPending},
	}

	t.Run("should always call repository when cache is disabled", func(t *testing.T) {
		mockRepo.EXPECT().
			ListByUserID(gomock.Any(), userID).
			Return(expectedTasks, nil)

		tasks, err := service.ListTasks(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
	})
}

func TestTaskService_ListTasks_CacheErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Get the address before closing
	addr := mr.Addr()
	mr.Close() // Close immediately to simulate connection failure

	// Create Redis client that will fail to connect
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	defer rdb.Close()

	userID := uint(1)
	expectedTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: model.TaskStatusPending},
	}

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	t.Run("should fallback to DB when cache fails", func(t *testing.T) {
		// Should call repository when cache fails
		mockRepo.EXPECT().
			ListByUserID(gomock.Any(), userID).
			Return(expectedTasks, nil)

		tasks, err := service.ListTasks(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
	})
}

func TestTaskService_ListTasks_Singleflight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	userID := uint(1)
	expectedTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task 1", Status: model.TaskStatusPending},
	}

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	t.Run("concurrent cache misses should only call repository once", func(t *testing.T) {
		// Expect repository to be called only once for concurrent requests
		mockRepo.EXPECT().
			ListByUserID(gomock.Any(), userID).
			Return(expectedTasks, nil).
			Times(1) // Should only be called once

		// Make concurrent requests
		done := make(chan bool, 3)
		for i := 0; i < 3; i++ {
			go func() {
				tasks, err := service.ListTasks(context.Background(), userID)
				require.NoError(t, err)
				assert.Equal(t, expectedTasks, tasks)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}
	})
}

func TestTaskService_CreateTask_CacheInvalidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	userID := uint(1)
	newTask := &model.Task{
		UserID: userID,
		Title:  "New Task",
		Status: model.TaskStatusPending,
	}

	// Pre-populate cache
	key := cache.KeyUserTasks(userID)
	existingTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Existing Task", Status: model.TaskStatusDone},
	}
	cachedData, err := json.Marshal(existingTasks)
	require.NoError(t, err)
	err = rdb.Set(context.Background(), key, cachedData, 60*time.Second).Err()
	require.NoError(t, err)

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	// Expect repository call for task creation
	mockRepo.EXPECT().
		CreateTask(gomock.Any(), newTask).
		Return(nil)

	// Create task
	err = service.CreateTask(context.Background(), newTask)
	require.NoError(t, err)

	// Verify cache was invalidated
	_, err = rdb.Get(context.Background(), key).Result()
	assert.Error(t, err) // Should return error as key should be deleted
}

func TestTaskService_UpdateTask_CacheInvalidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	userID := uint(1)
	updatedTask := &model.Task{
		ID:     1,
		UserID: userID,
		Title:  "Updated Task",
		Status: model.TaskStatusDone,
	}

	// Pre-populate cache
	key := cache.KeyUserTasks(userID)
	existingTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Original Task", Status: model.TaskStatusPending},
	}
	cachedData, err := json.Marshal(existingTasks)
	require.NoError(t, err)
	err = rdb.Set(context.Background(), key, cachedData, 60*time.Second).Err()
	require.NoError(t, err)

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	// Expect repository call for task update
	mockRepo.EXPECT().
		UpdateTask(gomock.Any(), updatedTask).
		Return(nil)

	// Update task
	err = service.UpdateTask(context.Background(), updatedTask)
	require.NoError(t, err)

	// Verify cache was invalidated
	_, err = rdb.Get(context.Background(), key).Result()
	assert.Error(t, err) // Should return error as key should be deleted
}

func TestTaskService_DeleteTask_CacheInvalidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Create a miniredis mock server
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Create Redis client connected to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	userID := uint(1)
	taskID := uint(1)

	// Pre-populate cache
	key := cache.KeyUserTasks(userID)
	existingTasks := []*model.Task{
		{ID: 1, UserID: userID, Title: "Task to Delete", Status: model.TaskStatusPending},
		{ID: 2, UserID: userID, Title: "Another Task", Status: model.TaskStatusDone},
	}
	cachedData, err := json.Marshal(existingTasks)
	require.NoError(t, err)
	err = rdb.Set(context.Background(), key, cachedData, 60*time.Second).Err()
	require.NoError(t, err)

	service := service.NewTaskService(mockRepo, rdb, 60*time.Second)

	// Expect repository calls for task deletion
	existingTask := &model.Task{ID: taskID, UserID: userID, Title: "Task to Delete"}
	mockRepo.EXPECT().
		FindByID(gomock.Any(), taskID).
		Return(existingTask, nil)
	mockRepo.EXPECT().
		DeleteTask(gomock.Any(), taskID).
		Return(nil)

	// Delete task
	err = service.DeleteTask(context.Background(), userID, taskID)
	require.NoError(t, err)

	// Verify cache was invalidated
	_, err = rdb.Get(context.Background(), key).Result()
	assert.Error(t, err) // Should return error as key should be deleted
}

func TestTaskService_CacheInvalidation_WhenCacheDisabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTaskRepository(ctrl)

	// Service without Redis (cache disabled)
	service := service.NewTaskService(mockRepo, nil, 60*time.Second)

	userID := uint(1)
	newTask := &model.Task{
		UserID: userID,
		Title:  "New Task",
		Status: model.TaskStatusPending,
	}

	// Expect repository call for task creation
	mockRepo.EXPECT().
		CreateTask(gomock.Any(), newTask).
		Return(nil)

	// Should not panic when cache is disabled
	err := service.CreateTask(context.Background(), newTask)
	require.NoError(t, err)
}
