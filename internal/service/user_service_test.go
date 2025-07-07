package service

import (
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().FindByEmail("test@example.com").Return(nil, nil)

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	svc := NewUserService(mockRepo)

	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "plaintextpassword",
	}

	err := svc.CreateUser(user)
	assert.Nil(t, err)
	assert.NotEqual(t, "plaintextpassword", user.PasswordHash)
}
