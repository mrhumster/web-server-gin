package repository

import (
	"context"

	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/sony/gobreaker"
	"gorm.io/gorm"
	"log"
	"time"
)

func ptrStr(s string) *string {
	return &s
}

var dbCircuitBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
	Name:        "Database",
	MaxRequests: 3,
	Interval:    30 * time.Second,
	Timeout:     60 * time.Second,
	ReadyToTrip: func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures > 3
	},
	OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
		log.Printf("Circuit Breaker %s: %s -> %s", name, from, to)
	},
})

func WithCircuitBreaker(fn func() error) error {
	_, err := dbCircuitBreaker.Execute(func() (interface{}, error) {
		return nil, fn()
	})
	return err
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (uint, error) {
	if user.Role == nil {
		user.Role = ptrStr("member")
	}
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (r *UserRepository) ReadUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user *models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uint, user request.UpdateUserRequest) (uint, error) {
	var userForUpdate *models.User
	result := r.db.WithContext(ctx).First(&userForUpdate, id)
	if result.Error != nil {
		return 0, result.Error
	}
	userForUpdate.FillInTheUpdateRequest(user)
	result = r.db.WithContext(ctx).Save(&userForUpdate)
	if result.Error != nil {
		return 0, result.Error
	}
	return userForUpdate.ID, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	return result.Error
}

func (r *UserRepository) ReadUserList(ctx context.Context, l, page int64) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	limit := int(l)
	offset := int((page - 1) * l)
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users)
	if result.Error != nil {
		return []models.User{}, int64(0), result.Error
	}
	return users, total, nil
}

func (r *UserRepository) ReadUserByEmail(ctx context.Context, value string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Model(&models.User{}).First(&user, "email = ?", value).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Exists(ctx context.Context, id uint) bool {
	u, _ := r.ReadUserByID(ctx, id)
	return u != nil
}
