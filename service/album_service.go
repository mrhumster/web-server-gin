package service

import (
	"context"

	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
)

type AlbumService struct {
	repo *repository.AlbumRepository
}

func NewAlbumService(repo *repository.AlbumRepository) *AlbumService {
	return &AlbumService{repo: repo}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, album models.Album) (uint, error) {
	return s.repo.CreateAlbum(ctx, album)
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, id uint) (*models.Album, error) {
	return s.repo.GetAlbumByID(ctx, id)
}

func (s *AlbumService) DeleteAlbumByID(ctx context.Context, id uint) error {
	return s.repo.DeleteAlbumByID(ctx, id)
}
