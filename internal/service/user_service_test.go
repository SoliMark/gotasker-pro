package service

import (
	"context"
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository/mock_repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	ctx := context.Background()

	mockRepo.EXPECT().
		FindByEmail(ctx, "test@example.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	svc := NewUserService(mockRepo)

	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "plaintextpassword",
	}

	err := svc.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotEqual(t, "plaintextpassword", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)
}
