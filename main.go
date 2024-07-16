package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var ENV_PORT = "PORT"

var db *sql.DB

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/api/secret", RevealSecretHandler).Methods("GET")

	r.HandleFunc("/api/user/sign-up", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/user/login", LoginHandler).Methods("POST")

	port := EnvVar(ENV_PORT)

	slog.Info(fmt.Sprintf("starting http server on port :%s", port))

	if err := os.Mkdir("db", os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Errorf("error creating the **db** dir : %v", err))
	}

	var err error

	db, err = sql.Open("sqlite3", "db/nandi.db")
	if err != nil {
		panic(fmt.Errorf("couldn't connect to db : %v", err))
	}
	defer db.Close()

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
func EnvVar(name string) string {
	return os.Getenv(name)
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

func errorResponse(code int, err string) Response {
	return Response{
		StatusCode: code,
		Status:     http.StatusText(code),
		Error:      err,
	}
}

func success(msg string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    msg,
	}
}
