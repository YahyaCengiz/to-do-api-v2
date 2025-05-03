package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/YahyaCengiz/todo-v2/middleware"
	"github.com/YahyaCengiz/todo-v2/models"
	"github.com/YahyaCengiz/todo-v2/store"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store := r.Context().Value("store").(*store.Store)

	var user *models.User
	for _, u := range store.Users {
		if u.Username == loginReq.Username && u.Password == loginReq.Password {
			user = &u
			break
		}
	}

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
