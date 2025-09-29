package service

import (
	"context"

	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (uint, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) ReadUser(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.ReadUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, user models.User) (uint, error) {
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUserByID(ctx, id)
}

func (s *UserService) ReadUserList(ctx context.Context) ([]models.User, error) {
	return s.repo.ReadUserList(ctx)
}
