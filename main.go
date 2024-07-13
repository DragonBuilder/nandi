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

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
// 	w.WriteHeader(e.httpStatus)
// 	return nil
// }

func ping(w http.ResponseWriter, r *http.Request) {
}

// func SignUp(w http.ResponseWriter, r *http.Request) {
// 	u := UserSignUp{}
// 	if err := render.Decode(r, &u); err != nil {
// 		slog.Error(fmt.Sprintf("ErrMashallingSignUpUser : %v", err))
// 		render.Render(w, r, ErrorResponse{"sign-up failed", http.StatusInternalServerError})
// 	}

// }

func EnvVar(name string) string {
	return os.Getenv(name)
}
