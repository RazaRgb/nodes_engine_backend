package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"userId"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashedPassword"`
	CreationDate   time.Time `json:"creationDate"`
}

type Project struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Owner        uuid.UUID `json:"owner"`
	CreationDate time.Time `json:"creationDate"`
	LastModified time.Time `json:"lastModified"`
}

// -----------------------------------------

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Data struct {
	Label string `json:"label"`
	Value string `json:"value"` //stored in db as json
}

// -----------------------------------------

type Tree struct {
	ID    string `json:"id"`
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID   string   `json:"id"`
	Type string   `json:"type"`
	Pos  Position `json:"position"`
	Data Data     `json:"data"`
}

type Edge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle,omitempty"`
	Target       string `json:"target"`
	TargetHandle string `json:"targetHandle,omitempty"`
}
