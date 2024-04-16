package api

import "net/http"

func (a *Api) Router() *http.ServeMux {
	router := http.NewServeMux()
	globalMiddlewares := a.chainMiddleware(a.loggingMiddleware, a.requestIDMiddleware)

	router.HandleFunc("POST /login", a.loginHandler)
	router.HandleFunc("POST /register", a.registerHandler)

	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", globalMiddlewares(router)))
	return v1Router
}
