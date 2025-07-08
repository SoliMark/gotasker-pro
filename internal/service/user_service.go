package service

import (
	"context"
	"errors"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	existing, _ := s.repo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := util.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword
	return s.repo.Create(ctx, user)
}
