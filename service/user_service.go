package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFount      = errors.New("user not found")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (uint, error) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, ErrUserAlreadyExists
			}
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (s *UserService) ReadUser(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.ReadUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, user request.UpdateUserRequest) (uint, error) {
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUserByID(ctx, id)
}

func (s *UserService) ReadUserList(ctx context.Context, limit, page int64) ([]models.User, int64, error) {
	return s.repo.ReadUserList(ctx, limit, page)
}
