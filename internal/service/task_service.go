package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"

	"github.com/SoliMark/gotasker-pro/internal/cache"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTask(ctx context.Context, id uint) (*model.Task, error)
	ListTasks(ctx context.Context, userID uint) ([]*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, userID, taskID uint) error
}

type taskService struct {
	repo    repository.TaskRepository
	rdb     *redis.Client
	ttl     time.Duration
	sfGroup singleflight.Group
}

func NewTaskService(repo repository.TaskRepository, rdb *redis.Client, ttl time.Duration) TaskService {
	return &taskService{
		repo:    repo,
		rdb:     rdb,
		ttl:     ttl,
		sfGroup: singleflight.Group{},
	}
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
	// fallback when cache is not enabled
	if s.rdb == nil {
		return s.repo.ListByUserID(ctx, userID)
	}
	key := cache.KeyUserTasks(userID)

	// fast path: cache hit
	if b, err := s.rdb.Get(ctx, key).Bytes(); err == nil && len(b) > 0 {
		var tasks []*model.Task
		if json.Unmarshal(b, &tasks) == nil {
			return tasks, nil
		}
	}

	// collapse concurrent misses
	v, err, _ := s.sfGroup.Do(key, func() (interface{}, error) {
		// double-check after acquiring singleflight
		if b, err := s.rdb.Get(ctx, key).Bytes(); err == nil && len(b) > 0 {
			var tasks []*model.Task
			if json.Unmarshal(b, &tasks) == nil {
				return tasks, nil
			}
		}
		// load from DB
		list, err := s.repo.ListByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		// set cache with TTL jitter (±10%)
		if data, e := json.Marshal(list); e == nil {
			jitter := cache.Jitter{}
			_ = s.rdb.Set(ctx, key, data, jitter.TTL(s.ttl, 0.1)).Err()
		}
		return list, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]*model.Task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, task *model.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("title is required")
	}
	return s.repo.UpdateTask(ctx, task)
}

func (s *taskService) DeleteTask(ctx context.Context, userID, taskID uint) error {
	t, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrTaskNotFound
	}
	if t.UserID != userID {
		return ErrPermissionDenied
	}
	return s.repo.DeleteTask(ctx, taskID)
}
