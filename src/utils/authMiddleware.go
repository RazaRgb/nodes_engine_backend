package utils

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userEmailKey contextKey = "email"

func AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized session", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWTSECRET")
		if secret == "" {
			log.Fatal("FATAL: jwt secret not set")
		}

		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized session", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized session, invalid token data", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Unauthorized session, missing email", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userEmailKey, email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}
