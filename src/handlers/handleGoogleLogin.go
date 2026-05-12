package handlers

import (
	//"backend/src/db"
	//"backend/src/models"
	"backend/src/engine"
	"backend/src/utils"
	"fmt"
	"time"

	//"encoding/json"
	//"fmt"
	//"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	"net/http"
)

func HandleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	stateCookie := r.Cookie("oauth_state")

	utils.JsonWrite(w, nil, http.StatusOK)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		fmt.Printf("unable to generate random string for google login")
		http.Error(w, "unable to login", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 10),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	url := engine.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
