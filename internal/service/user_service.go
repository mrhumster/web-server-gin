package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrhumster/web-server-gin/internal/delivery/http/dto/request"
	"github.com/mrhumster/web-server-gin/internal/domain/models"
	"github.com/mrhumster/web-server-gin/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFount      = errors.New("user not found")
)

type UserService struct {
	repo              repository.UserRepository
	permissionService *PermissionService
	mu                sync.RWMutex
}

func NewUserService(repo repository.UserRepository, perm *PermissionService) *UserService {
	return &UserService{repo: repo, permissionService: perm}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	if user.Role == "" {
		role := "member"
		user.Role = role
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, ErrUserAlreadyExists
			}
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	s.mu.Lock()
	policy := id.String()
	resource := fmt.Sprintf("users/%s", id.String())
	log.Printf("⚠️ UserService. CreateUser Permission debug: %s %s", policy, resource)
	s.permissionService.AddPolicy(policy, resource, "read")
	s.permissionService.AddPolicy(policy, resource, "write")
	s.permissionService.AddPolicy(policy, resource, "delete")
	s.permissionService.AddPolicy(policy, "users", "read")
	s.permissionService.AddPolicy(policy, "stream", "read")
	s.permissionService.AddPolicy(policy, "stream", "write")
	s.mu.Unlock()
	return id, nil
}

func (s *UserService) ReadUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.ReadUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, user request.UpdateUserRequest) (*uuid.UUID, error) {
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteUserByID(ctx, id)
	if err == nil {
		policy := fmt.Sprintf("%s", id)
		resource := fmt.Sprintf("users/%s", id.String())
		s.permissionService.RemovePolicy(policy, resource, "read")
		s.permissionService.RemovePolicy(policy, resource, "write")
		s.permissionService.RemovePolicy(policy, resource, "delete")
		s.permissionService.RemovePolicy(policy, "users", "read")
		s.permissionService.RemovePolicy(policy, "stream", "read")
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

func (s *UserService) UpdateTokenVersion(ctx context.Context, userID *uuid.UUID, version string) error {
	return s.repo.UpdateTokenVersion(ctx, userID, version)
}
