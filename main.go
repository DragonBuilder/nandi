package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var ENV_PORT = "PORT"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", ping)
	})

	port := EnvVar(ENV_PORT)

	slog.Info(fmt.Sprintf("starting http server on port :%s", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"message": "pong"})
}

func EnvVar(name string) string {
	return os.Getenv(name)
}
