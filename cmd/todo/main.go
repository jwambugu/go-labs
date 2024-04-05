package main

import (
	"context"
	"encoding/json"
	httpapi "go-labs/internal/http"
	"go-labs/internal/repository/db"
	"go-labs/svc/auth"
	"log"
	"net/http"
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

type Middleware func(http.Handler) http.Handler

func CreateMiddleware(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}

		return next
	}
}

type writer struct {
	http.ResponseWriter
	statusCode int
}

func (w *writer) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &writer{
			ResponseWriter: w,
			statusCode:     200,
		}

		defer func() {
			log.Printf("[%d] %s %s %s", writer.statusCode, r.Method, r.URL.Path, time.Since(start))
		}()

		next.ServeHTTP(writer, r)
	})
}

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Allowing CORS")
		next.ServeHTTP(w, r)
	})
}

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking IsAuth")
		next.ServeHTTP(w, r)
	})
}

func CheckPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking Permission")
		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}

	//middlewares := CreateMiddleware(
	//	Logger,
	//	AllowCors,
	//	IsAuth,
	//	CheckPermission,
	//)

	router := http.NewServeMux()

	router.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var loginRequest *auth.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			httpapi.NewError(err, http.StatusBadRequest).Send(w)
			return
		}

		log.Printf("login: %#v", loginRequest)

		if err := loginRequest.Validate(); err != nil {
			httpapi.NewError(err, http.StatusBadRequest).Send(w)
			return
		}

	})

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	srv := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("http: %v", err)
	}
}
