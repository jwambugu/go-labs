package db

import (
	"context"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"go-labs/internal/config"
	"path/filepath"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate runs the migrations
func Migrate(db *sqlx.DB, driver string) error {
	goose.SetBaseFS(migrationFiles)

	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("set dialect: %v", err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %v", err)
	}

	return nil
}

func New(ctx context.Context, driver string, dbName string) (*sqlx.DB, error) {
	pwd, err := config.GetPWD()
	if err != nil {
		return nil, err
	}

	dsn := filepath.Join("file:", pwd, dbName)
	dsn = dsn + `?cache=shared&mode=rwc`

	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("db connect: %v", err)
	}

	return db, nil
}
