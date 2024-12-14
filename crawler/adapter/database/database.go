package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-manager/vision-seeker/crawler/domain/image"
)

type Database interface {
	Close()
	InsertUrl(hash string) (bool, error)
	InsertImage(hash string, img *image.Image) (bool, error)
	InsertLabel(hash, label string) (bool, error)
	InsertMapping(imgHash, lblHash string) (bool, error)
}

type database struct {
	conn *pgxpool.Pool
}

func New(host, port, name, user, pass string) (*database, error) {
	conn, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name),
	)
	if err != nil {
		return nil, err
	}
	return &database{conn: conn}, nil
}

func (db *database) Close() {
	db.conn.Close()
}

func insertResult(err error) (bool, error) {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// unique_violation
			if pgErr.Code == "23505" {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (db *database) InsertUrl(hash string) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO "visited" (hash) VALUES ($1);`,
		hash,
	)
	return insertResult(err)
}

func (db *database) InsertImage(hash string, img *image.Image) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO "image" (hash, size, width, height, entropy, format) 
			VALUES ($1, $2, $3, $4, $5, $6);`,
		hash,
		img.Size,
		img.Width,
		img.Height,
		img.Entropy(),
		img.Format,
	)
	return insertResult(err)
}

func (db *database) InsertLabel(hash, label string) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO "label" (hash, label) VALUES ($1, $2);`,
		hash,
		label,
	)
	return insertResult(err)
}

func (db *database) InsertMapping(imgHash, lblHash string) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO "image_label_mapping" 
			(image_hash, label_hash) VALUES ($1, $2);`,
		imgHash,
		lblHash,
	)
	return insertResult(err)
}
