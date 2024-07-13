package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

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

	if len(strings.TrimSpace(signUp.Email)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Email not provided"})
		slog.Error(fmt.Sprintf("Email not provided : %v", err))
		return
	}

	if len(strings.TrimSpace(signUp.Password)) == 0 || len(strings.TrimSpace(signUp.ConfirmPassword)) == 0 {
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

	json.NewEncoder(w).Encode(map[string]string{"Message": fmt.Sprintf("Registered user %s", signUp.Email)})
}
