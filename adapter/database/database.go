package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Close()
	InsertUrl(hash string) (bool, error)
}

type database struct {
	conn *pgxpool.Pool
}

func (db *database) createTables() error {
	_, err := db.conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS "url" (
		hash VARCHAR(64) PRIMARY KEY
	);`)
	return err
}

func New(host, port, name, user, pass string) (*database, error) {
	conn, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name),
	)
	if err != nil {
		return nil, err
	}
	db := &database{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *database) Close() {
	db.conn.Close()
}

func (db *database) InsertUrl(hash string) (bool, error) {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO "url" (hash) VALUES ($1);`,
		hash,
	)
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
