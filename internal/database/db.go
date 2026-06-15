package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	return sql.Open(
		"sqlite3",
		"./data/cli-auth.db",
	)
}
