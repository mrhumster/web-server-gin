package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrhumster/web-server-gin/models"
)

type AlbumRepository struct {
	db *pgxpool.Pool
}

func NewAlbumRepository(db *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album models.Album) (int, error) {
	var id int
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()
	query := `INSERT INTO albums (title, artist, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query, album.Title, album.Artist, album.Price, album.CreatedAt, album.UpdatedAt).Scan(&id)
	return id, err
}
