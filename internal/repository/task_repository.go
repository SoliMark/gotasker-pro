package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/SoliMark/gotasker-pro/internal/model"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id uint) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id uint) error
	ListByUserID(ctx context.Context, userID uint) ([]*model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) FindByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}

func (r *taskRepository) ListByUserID(ctx context.Context, userID uint) ([]*model.Task, error) {
	var tasks []*model.Task
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}
