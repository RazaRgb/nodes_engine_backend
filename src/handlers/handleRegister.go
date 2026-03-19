package handlers

import (
	// "fmt"
	"fmt"
	"net/http"

	// "time"
	"backend/src/db"
	"backend/src/models"
	"backend/src/utils"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	//validation 
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(username) < 5 || len(password) < 8 {
		http.Error(w, "Invalid username or password", http.StatusNotAcceptable)
		return
	}

	userExists, _ := db.UserExists(email)
	if userExists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.PassHash(password)
	if err != nil {
		fmt.Printf("password hashing failed while registering user %v", err)
		http.Error(w, "unable to create user", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
	}

	err = db.InsertUser(user)
	if err != nil {
		fmt.Printf("user insertion in DB failed %v", err)
		http.Error(w, "unable to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"User successfully registered"}`))
}
