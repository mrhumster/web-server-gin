package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

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

func strPtr(s string) *string {
	return &s
}

type UserService struct {
	repo              *repository.UserRepository
	permissionService *PermissionService
	mu                sync.RWMutex
}

func NewUserService(repo *repository.UserRepository, perm *PermissionService) *UserService {
	return &UserService{repo: repo, permissionService: perm}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (uint, error) {
	if user.Role == nil {
		role := "member"
		user.Role = &role
	}
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

	s.mu.Lock()
	policy := fmt.Sprintf("%d", id)
	resource := fmt.Sprintf("users/%d", id)
	s.permissionService.AddPolicy(policy, resource, "*")
	s.permissionService.AddPolicy(policy, "users", "read")
	s.mu.Unlock()
	return id, nil
}

func (s *UserService) ReadUser(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.ReadUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, user request.UpdateUserRequest) (uint, error) {
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.DeleteUserByID(ctx, id)
	if err == nil {
		policy := fmt.Sprintf("users:%d", id)
		resource := fmt.Sprintf("users/%d", id)
		s.permissionService.RemovePolicy(policy, resource, "*")
	}
	return err

}

func (s *UserService) ReadUserList(ctx context.Context, limit, page int64) ([]models.User, int64, error) {
	return s.repo.ReadUserList(ctx, limit, page)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.ReadUserByEmail(ctx, email)
}

func (s *UserService) ValidateUser(ctx context.Context, email, password string) (*models.User, error) {
	var user *models.User
	var err error
	if user, err = s.GetUserByEmail(ctx, email); err != nil {
		return nil, err
	}

	if user.CheckPassword(password) {
		return user, nil
	}
	return nil, errors.New("invalid password")
}

func (s *UserService) UpdateTokenVersion(ctx context.Context, userID uint64, version string) error {
	return s.repo.UpdateTokenVersion(ctx, userID, version)
}
