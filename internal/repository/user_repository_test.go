package repository_test

import (
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test DB: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("fail to migrate:%v", err)
	}
	return db
}

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &model.User{
		Username:     "sol",
		Email:        "sol@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)

	found, err := repo.FindByEmail("sol@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "sol", found.Username)
	assert.True(t, user.CreatedAt.Equal(found.CreatedAt))
	assert.True(t, user.UpdatedAt.Equal(found.UpdatedAt))

	foundByID, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "sol", foundByID.Username)
}
