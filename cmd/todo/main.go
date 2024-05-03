package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-labs/internal/api"
	"go-labs/internal/repository"
	"go-labs/internal/repository/db"
	"go-labs/svc/auth"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), interruptSignals...)
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
		logger.Fatal("failed to create jwt manager", zap.Error(err))
	}

	repoStore := repository.NewStore()
	repoStore.User = repository.NewUserRepo(app.db)

	authSVC := auth.NewAuthSvc(logger, repoStore, jwtManager)

	httpApi := api.NewApi(repoStore, authSVC)
	srv := httpApi.Serve(8080)

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		if err = srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return fmt.Errorf("start http server: %v", err)
		}
		return nil
	})

	errGroup.Go(func() error {
		<-ctx.Done()

		logger.Info("http server shutdown")

		if err = srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("http shutdown: %v", err)
		}

		logger.Info("http server shutdown successfully")

		return nil
	})

	if err = errGroup.Wait(); err != nil {
		logger.Fatal("error group", zap.Error(err))
	}
}
