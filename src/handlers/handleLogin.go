package handlers

import (
	"backend/src/db"
	"backend/src/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := db.GetUser(email)
	if err != nil {
		fmt.Printf("error getting user %+v\n", err)
		http.Error(w, "failed to log in", http.StatusInternalServerError)
		return
	}
	if !utils.CheckPassHash(password, user.HashedPassword) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.GenerateToken(email)
	if err != nil {
		http.Error(w, "unable to log in", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		//SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}
