package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://postgres:lol@localhost:5432/gosql")

	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Database is connected.....")

	pingErr := db.Ping()

	if pingErr != nil {
		return nil, fmt.Errorf("db: open %w", pingErr)
	}

	return db, nil

}

func MigrateFS(db *sql.DB, baseFS fs.FS, dir string) error {
	goose.SetBaseFS(baseFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("db: migrate %w", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("db: goose up %w", err)
	}

	return nil
}
