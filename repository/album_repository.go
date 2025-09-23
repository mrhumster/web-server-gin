package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mrhumster/web-server-gin/models"
)

type AlbumRepository struct {
	db *pgx.Conn
}

func NewAlbumRepository(db *pgx.Conn) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album models.Album) (int, error) {
	var id int
	query := `INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, album.Title, album.Artist, album.Price).Scan(&id)
	return id, err
}
