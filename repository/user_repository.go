package repository

import (
	"context"

	"github.com/mrhumster/web-server-gin/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (uint, error) {
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

func (r *UserRepository) UpdateUser(ctx context.Context, id uint, user models.User) (uint, error) {
	var userForUpdate *models.User
	result := r.db.WithContext(ctx).First(&userForUpdate, id)
	if result.Error != nil {
		return 0, result.Error
	}

	userForUpdate.Email = user.Email
	userForUpdate.Login = user.Login

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
