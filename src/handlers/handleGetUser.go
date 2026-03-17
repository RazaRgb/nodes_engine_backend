package handlers

import (
	"backend/src/db"
	"backend/src/utils"
	"encoding/json"
	"net/http"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	email, ok := utils.GetEmailFromContext(r.Context())
	if !ok {
		http.Error(w, "user doesnt exist", http.StatusNotFound)
		return
	}
	user, err := db.GetUser(email)
	if err != nil {
		http.Error(w, "user doesnt exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"username":        user.Username,
		"email":           user.Email,
		"account_created": user.CreationDate,
	})
}
