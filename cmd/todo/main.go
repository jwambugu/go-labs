package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go-labs/internal/api"
	"go-labs/internal/repository"
	"go-labs/internal/repository/db"
	"go-labs/svc/auth"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

type app struct {
	db *sqlx.DB
}

func run(ctx context.Context) (*app, error) {
	dbConn, err := db.New(ctx, "sqlite3", "todo.db")
	if err != nil {
		return nil, err
	}

	if err = db.Migrate(dbConn, "sqlite3"); err != nil {
		return nil, err
	}

	newApp := &app{
		db: dbConn,
	}

	return newApp, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	app, err := run(ctx)
	if err != nil {
		log.Fatalf("run: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// TODO: read from the env
	pasetoKey := "86f5778df1b11e35caf8bc793391bfd1"

	jwtManager, err := auth.NewPasetoToken(pasetoKey)
	if err != nil {
		logger.Error("failed to create jwt manager", zap.Error(err))
		log.Panicln(err)
	}

	repoStore := repository.NewStore()
	repoStore.User = repository.NewUserRepo(app.db)

	authSVC := auth.NewAuthSvc(logger, repoStore, jwtManager)

	httpApi := api.NewApi(repoStore, authSVC)
	srv := httpApi.Serve(8080)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("start server: %v", err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
}
