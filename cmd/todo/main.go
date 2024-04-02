package main

import (
	"context"
	"go-labs/internal/repository/db"
	"log"
	"time"
)

func run(ctx context.Context) error {
	dbConn, err := db.New(ctx, "sqlite3", "todo.db")
	if err != nil {
		return err
	}

	if err = db.Migrate(dbConn, "sqlite3"); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
