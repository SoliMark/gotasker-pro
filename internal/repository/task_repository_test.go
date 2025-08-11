package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
)

// 使用 SQLite 的 in-memory DB 初始化
func setupSQLiteTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自動 migrate Task model
	err = db.AutoMigrate(&model.Task{})
	assert.NoError(t, err)

	return db
}

func TestTaskRepository_CRUD_SQLite(t *testing.T) {
	db := setupSQLiteTestDB(t)
	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	// Create
	task := &model.Task{
		UserID:  42,
		Title:   "SQLite Task",
		Content: "This is a task in SQLite",
		Status:  "pending",
	}
	err := repo.CreateTask(ctx, task)
	assert.NoError(t, err)
	assert.NotZero(t, task.ID)

	// FindByID
	found, err := repo.FindByID(ctx, task.ID)
	assert.NoError(t, err)
	assert.Equal(t, task.Title, found.Title)

	// Update
	found.Title = "Updated Title"
	err = repo.UpdateTask(ctx, found)
	assert.NoError(t, err)

	updated, err := repo.FindByID(ctx, task.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)

	// ListByUserID
	list, err := repo.ListByUserID(ctx, task.UserID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// Delete
	err = repo.DeleteTask(ctx, task.ID)
	assert.NoError(t, err)

	deleted, err := repo.FindByID(ctx, task.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
