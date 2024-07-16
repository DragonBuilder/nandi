package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		slog.Error(fmt.Sprintf("ErrTemplateExecute : %v", err))
		return
	}
}

type UserSignUp struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error reading sign-up request body : %v", err))
		return
	}

	var signUp UserSignUp
	if err := json.Unmarshal(body, &signUp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error unmarshalling sign-up request body : %v", err))
		return
	}

	signUp = UserSignUp{
		Email:           strings.TrimSpace(signUp.Email),
		Password:        strings.TrimSpace(signUp.Password),
		ConfirmPassword: strings.TrimSpace(signUp.ConfirmPassword),
	}

	if len(signUp.Email) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Email not provided"})
		slog.Error(fmt.Sprintf("Email not provided : %v", err))
		return
	}

	if len(signUp.Password) == 0 || len(signUp.ConfirmPassword) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Provide password and the confirm password. Both should be same."})
		slog.Error(fmt.Sprintf("Email not provided : %v", err))
		return
	}

	hPwd, err := bcrypt.GenerateFromPassword([]byte(signUp.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error generating password hash : %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword(hPwd, []byte(signUp.ConfirmPassword)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Password mismatch. Make sure password & confirm password are same."})
		slog.Error(fmt.Sprintf("Error password mismatch : %v", err))
		return
	}

	db, err := sql.Open("sqlite3", "db/nandi.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error"})
		slog.Error(fmt.Sprintf("Error opening connection to db : %v", err))
		return
	}
	defer db.Close()

	const create = `
  CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
  );`

	if _, err := db.Exec(create); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error"})
		slog.Error(fmt.Sprintf("Error creating users table : %v", err))
		return
	}

	now := time.Now()
	if _, err := db.Exec("INSERT INTO users VALUES(NULL, ?, ?, ?, ?)", signUp.Email, string(hPwd), now, now); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error"})
		slog.Error(fmt.Sprintf("Error inserting users table : %v", err))
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"Message": fmt.Sprintf("Registered user %s", signUp.Email)})
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp := errorResponse(http.StatusBadRequest, "error reading body")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(fmt.Sprintf("%s : %v", resp.Error, err))
		return
	}

	var login Login
	if err := json.Unmarshal(body, &login); err != nil {
		var resp = errorResponse(http.StatusBadRequest, "error bad type")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(fmt.Sprintf("%s : %v", resp.Error, err))
		return
	}

	login = Login{
		Email:    strings.TrimSpace(login.Email),
		Password: strings.TrimSpace(login.Password),
	}

	fetch := `
SELECT password FROM users WHERE email=?
`
	rows, err := db.Query(fetch, login.Email)
	if err != nil {
		var resp = errorResponse(http.StatusBadRequest, fmt.Sprintf("unknown user : %s", login.Email))
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(fmt.Sprintf("%s : %v", resp.Error, err))
		return
	}

	if !rows.Next() {
		var resp = errorResponse(http.StatusBadRequest, fmt.Sprintf("unknown user : %s", login.Email))
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(resp.Error)
		return
	}

	var pwd string
	if err := rows.Scan(&pwd); err != nil {
		var resp = errorResponse(http.StatusInternalServerError, "couldn't verify password")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(resp.Error)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(login.Password)); err != nil {
		var resp = errorResponse(http.StatusBadRequest, "bad credentials")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		slog.Error(resp.Error)
		return
	}

	json.NewEncoder(w).Encode(success("login successful"))

}
