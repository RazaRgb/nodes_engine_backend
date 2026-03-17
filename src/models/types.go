package models

import (
	"time"
)

type User struct {
	Email          string `json:"email"` // primary key
	Username       string `json:"username"`
	HashedPassword string `json:"hashedPassword"`
	// AuthProvider   string `json:"authProvider"`
	// AuthProviderID string `json:"authProviderID"`
	CreationDate time.Time `json:"creationDate"`
}

type Project struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Owner        string    `json:"owner"`
	CreationDate time.Time `json:"creationDate"`
	LastModified time.Time `json:"lastModified"`
}
