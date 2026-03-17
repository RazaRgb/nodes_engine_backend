package db

import (
	"backend/src/models"
	"context"
	"fmt"
	// "net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dbURL := os.Getenv("DBURL")

	// Create the connection pool
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database is unreachable: %v\n", err)
		os.Exit(1)
	}

	DB = pool
	fmt.Println("Successfully connected to Postgres 🚀")

	// TODO: add other tables
	query := `
	CREATE TABLE IF NOT EXISTS users (
		email VARCHAR(255) PRIMARY KEY,
		username VARCHAR(100) NOT NULL,
		hashedpass TEXT NOT NULL,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS projects (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL REFERENCES users(email) ON DELETE CASCADE,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
	`

	// Execute the query using our connection pool
	_, err = DB.Exec(context.Background(), query)
	if err != nil {
		fmt.Printf("Failed to create tables: %v\n", err)
		return
	}

	fmt.Println("Database schema initialized successfully!")
}

func UserExists(email string) (bool, error) {
	userQuery := `
	SELECT EXISTS (
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`
	row := DB.QueryRow(context.Background(), userQuery, email)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userExists query failed %w", err)
	}

	return exists, nil
}

func InsertUser(newUser models.User) error {
	insertUserQuery := `
	INSERT INTO users (email, username, hashedpass)
	VALUES ($1,$2,$3)
	`
	_, err := DB.Exec(
		context.Background(),
		insertUserQuery,
		newUser.Email,
		newUser.Username,
		newUser.HashedPassword,
	)
	if err != nil {
		fmt.Printf("failed to create user %s", newUser.Username)
		return fmt.Errorf("failed to create user %w", err)
	}

	return nil
}

func GetUser(email string) (models.User, error) {
	selectQuery := `
	SELECT email, username, hashedpass, creation_date
	FROM users
	WHERE email = $1
	`
	var user models.User
	row := DB.QueryRow(context.Background(), selectQuery, email)
	err := row.Scan(
		&user.Email,
		&user.Username,
		&user.HashedPassword,
		&user.CreationDate,
	)
	if err != nil {
		return user, fmt.Errorf("unable to get user %w", err)
	}

	return user, nil
}

func GetProjectsInfo(email string) ([]models.Project, error) {

	selectQuery := `
	SELECT name, creation_date, last_modified
	FROM projects
	WHERE owner = $1
	`

	rows, err := DB.Query(
		context.Background(),
		selectQuery,
		email,
	)
	if err != nil {
		fmt.Printf("error in getting projlist from db")
		return nil, err
	}
	defer rows.Close()

	projList := make([]models.Project, 0)

	for rows.Next() {

		var proj models.Project
		err := rows.Scan(
			&proj.Name,
			&proj.CreationDate,
			&proj.LastModified,
		)
		if err != nil {
			fmt.Printf("error scanning project list query: %+v\n", err)
			return nil, err
		}
		projList = append(projList, proj)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("error while row scan loop %+v\n", err)
		return nil, err
	}

	return projList, nil
}

func InsertProject(newProject models.Project) (models.Project, error) {
	insertQuery := `
	INSERT INTO projects (id, name, owner)
	values ($1,$2,$3)
	RETURNING creation_date, last_modified 
	`
	err := DB.QueryRow(
		context.Background(),
		insertQuery,
		newProject.ID,
		newProject.Name,
		newProject.Owner,
	).Scan(
		&newProject.CreationDate,
		&newProject.LastModified,
	)
	if err != nil {
		fmt.Printf("error while inserting project into db %+v \n", err)
		return newProject, err
	}
	return newProject , nil
}

func DeleteProject(projectID string, email string) error {
	deleteQuery :=	`
	DELETE FROM projects
	WHERE id = $1 AND owner = $2
	`
	result ,err := DB.Exec(
		context.Background(),
		deleteQuery,
		projectID,
		email,
		)
	if err != nil {
		fmt.Printf("unable to delete project %+v\n", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No project to delete / unauthorized deletion")
	}

	return nil
}
