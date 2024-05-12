package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(conn string) (*sql.DB, error) {
	const op = "db.New"

	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
