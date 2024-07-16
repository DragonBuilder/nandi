package main

import (
	"errors"
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
	r.HandleFunc("/api/user/login", LoginHandler).Methods("POST")

	port := EnvVar(ENV_PORT)

	slog.Info(fmt.Sprintf("starting http server on port :%s", port))

	if err := os.Mkdir("db", os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Errorf("error creating the **db** dir : %v", err))
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

func EnvVar(name string) string {
	return os.Getenv(name)
}
