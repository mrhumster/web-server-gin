package repository

import (
	"context"

	"github.com/mrhumster/web-server-gin/models"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album models.Album) (uint, error) {
	result := r.db.WithContext(ctx).Create(&album)
	if result.Error != nil {
		return 0, result.Error
	}
	return album.ID, nil
}

func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id uint) (*models.Album, error) {
	var album models.Album
	result := r.db.WithContext(ctx).First(&album, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &album, nil
}
