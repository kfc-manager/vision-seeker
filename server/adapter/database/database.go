package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-manager/vision-seeker/server/domain/image"
)

type Database interface {
	Close()
	InsertImage(img *image.Image) (bool, error)
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
	err = conn.Ping(context.Background())
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

func (db *database) InsertImage(img *image.Image) (bool, error) {
	label := pgtype.Text{}
	if len(img.Label) < 1 {
		label = pgtype.Text{Valid: false, String: ""}
	} else {
		label = pgtype.Text{Valid: true, String: img.Label}
	}

	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO image(hash, size, width, height, entropy, transparent, format, label) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		img.Hash, img.Size, img.Width, img.Height, img.Entropy(),
		img.Trans(), img.Format, label,
	)

	return insertResult(err)
}
