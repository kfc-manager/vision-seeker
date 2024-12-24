package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-manager/vision-seeker/storage/server/domain"
)

type Database interface {
	Close()
	InsertImgHash(i *domain.Image) (bool, error)
	ExistImgHash(hash string) (bool, error)
	InsertDataset(d *domain.Dataset) (bool, error)
	InsertImgMap(m *domain.Mapping) (bool, error)
	InsertLabel(l *domain.Label) (bool, error)
	InsertLabelMap(l *domain.Label) (bool, error)
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

func (db *database) InsertImgHash(i *domain.Image) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO image_hash(hash, created_at) VALUES ($1, $2);`,
		i.Hash, i.CreatedAt,
	)
	return insertResult(err)
}

func (db *database) ExistImgHash(hash string) (bool, error) {
	row := db.conn.QueryRow(
		context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM image_hash where hash = $1
		);`,
		hash,
	)
	exist := false
	if err := row.Scan(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

func (db *database) InsertDataset(d *domain.Dataset) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO dataset(id, name, created_at) VALUES ($1, $2, $3);`,
		d.Id, d.Name, d.CreatedAt,
	)
	return insertResult(err)
}

func (db *database) InsertImgMap(m *domain.Mapping) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO dataset_image_mapping(dataset_id, image_hash, created_at) 
			VALUES ($1, $2, $3);`,
		m.Id, m.Hash, m.CreatedAt,
	)
	return insertResult(err)
}

func (db *database) InsertLabel(l *domain.Label) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO label(hash, label) VALUES ($1, $2);`,
		l.Hash, l.Label,
	)
	return insertResult(err)
}

func (db *database) InsertLabelMap(l *domain.Label) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO dataset_image_label_mapping(dataset_id, image_hash, label_hash) 
			VALUES ($1, $2, $3);`,
		l.Mapping.Id, l.Mapping.Hash, l.Hash,
	)
	return insertResult(err)
}
