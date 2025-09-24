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

func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id int64) (*models.Album, error) {
	album := &models.Album{}
	query := `SELECT title, artist, price, created_at, updated_at FROM albums WHERE ID=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&album.Title,
		&album.Artist,
		&album.Price,
		&album.CreatedAt,
		&album.UpdatedAt)
	if err != nil {
		return nil, err
	}
	album.ID = id
	return album, err
}
