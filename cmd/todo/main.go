package main

import (
	"context"
	"go-labs/internal/api"
	"go-labs/internal/repository/db"
	"log"
	"os"
	"os/signal"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}

	httpApi := api.NewApi(nil)
	srv := httpApi.Serve(8080)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("start server: %v", err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
}
