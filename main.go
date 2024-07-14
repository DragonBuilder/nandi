package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var ENV_PORT = "PORT"

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	r.HandleFunc("/", HomeHandler).Methods("GET")

	r.HandleFunc("/api/user/sign-up", SignUpHandler).Methods("POST")

	port := EnvVar(ENV_PORT)

	slog.Info(fmt.Sprintf("starting http server on port :%s", port))

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func EnvVar(name string) string {
	return os.Getenv(name)
}
