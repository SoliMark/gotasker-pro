package service

import (
	"context"
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository/mock_repository"
	"github.com/SoliMark/gotasker-pro/internal/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	jwtMaker := util.NewJWTMaker("test_secret_key")

	ctx := context.Background()

	mockRepo.EXPECT().
		FindByEmail(ctx, "test@example.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	svc := NewUserService(mockRepo, jwtMaker)

	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "plaintextpassword",
	}

	err := svc.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotEqual(t, "plaintextpassword", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)
}

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	jwtMaker := util.NewJWTMaker("test_secret_key")

	svc := NewUserService(mockRepo, jwtMaker)
	ctx := context.Background()

	email := "test@example.com"
	password := "secret"
	hashedPassword, _ := util.HashPassword(password)

	user := &model.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil)

		token, err := svc.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByEmail(ctx, "notfound@example.com").
			Return(nil, nil)

		token, err := svc.AuthenticateUser(ctx, "notfound@example.com", password)
		require.ErrorIs(t, err, ErrUserNotFound)
		require.Empty(t, token)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		token, err := svc.AuthenticateUser(ctx, email, "wrongpassword")
		require.ErrorIs(t, err, ErrInvalidCredential)
		require.Empty(t, token)
	})
}
