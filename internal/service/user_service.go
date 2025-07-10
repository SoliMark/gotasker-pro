package service

import (
	"context"
	"errors"
	"time"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/repository"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invaild credentials")
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	AuthenticateUser(ctx context.Context, email, passord string) (string, error)
}

type userService struct {
	repo     repository.UserRepository
	jwtMaker *util.JWTMaker
}

func NewUserService(r repository.UserRepository, jwtMaker *util.JWTMaker) UserService {
	return &userService{
		repo:     r,
		jwtMaker: jwtMaker,
	}
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

func (s *userService) AuthenticateUser(ctx context.Context, email, passord string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrUserNotFound
	}

	if !util.CheckPasswordHash(passord, user.PasswordHash) {
		return "", ErrInvalidCredential
	}

	token, err := s.jwtMaker.GenerateToken(user.ID, time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}
