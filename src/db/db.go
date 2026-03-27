package db

import (
	"backend/src/models"
	"context"
	"fmt"
	// "net/http"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
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

	CREATE TABLE IF NOT EXISTS trees (
		id VARCHAR(255) PRIMARY KEY,
		project_id VARCHAR(255) NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS nodes (
		id VARCHAR(255) PRIMARY KEY,
		belongs_to VARCHAR(255) NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
		type VARCHAR(255) NOT NULL,
		pos_x REAL NOT NULL,
		pos_y REAL NOT NULL,
		label TEXT NOT NULL,
		value TEXT
	);

	CREATE TABLE IF NOT EXISTS edges (
		id VARCHAR(255) PRIMARY KEY,
		belongs_to VARCHAR(255) NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
		source VARCHAR(255) NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
		source_handle VARCHAR(255) NOT NULL,
		target VARCHAR(255) NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
		target_handle VARCHAR(255) NOT NULL
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

func UserExists(email string, tx ...pgx.Tx) (bool, error) {
	userQuery := `
	SELECT EXISTS (
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`
	var row pgx.Row
	if len(tx) == 0 {
		row = DB.QueryRow(context.Background(), userQuery, email)
	} else {
		row = tx[0].QueryRow(context.Background(), userQuery, email)
	}

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userExists query failed %w", err)
	}

	return exists, nil
}

func InsertUser(newUser models.User, tx ...pgx.Tx) error {
	insertUserQuery := `
	INSERT INTO users (email, username, hashedpass)
	VALUES ($1,$2,$3)
	`
	var err error
	if len(tx) == 0 {
		_, err = DB.Exec(
			context.Background(),
			insertUserQuery,
			newUser.Email,
			newUser.Username,
			newUser.HashedPassword,
		)
	} else {
		_, err = tx[0].Exec(
			context.Background(),
			insertUserQuery,
			newUser.Email,
			newUser.Username,
			newUser.HashedPassword,
		)
	}

	if err != nil {
		fmt.Printf("failed to create user %s\n", newUser.Username)
		return fmt.Errorf("failed to create user %w\n", err)
	}
	return nil
}

func GetUser(email string, tx ...pgx.Tx) (models.User, error) {
	selectQuery := `
	SELECT email, username, hashedpass, creation_date
	FROM users
	WHERE email = $1
	`
	var user models.User
	var row pgx.Row

	if len(tx) == 0 {
		row = DB.QueryRow(context.Background(), selectQuery, email)
	} else {
		row = tx[0].QueryRow(context.Background(), selectQuery, email)
	}
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

func GetProjectsInfo(email string, tx ...pgx.Tx) ([]models.Project, error) {
	selectQuery := `
	SELECT id, name, creation_date, last_modified
	FROM projects
	WHERE owner = $1
	`

	var rows pgx.Rows
	var err error
	if len(tx) == 0 {
		rows, err = DB.Query(context.Background(), selectQuery, email)
	} else {
		rows, err = tx[0].Query(context.Background(), selectQuery, email)
	}
	if err != nil {
		fmt.Printf("error in getting projlist from db")
		return nil, err
	}
	defer rows.Close()

	projList := make([]models.Project, 0)

	for rows.Next() {

		var proj models.Project
		err := rows.Scan(
			&proj.ID,
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

func InsertProject(newProject models.Project, tx ...pgx.Tx) (models.Project, error) {
	insertQuery := `
	INSERT INTO projects (id, name, owner)
	values ($1,$2,$3)
	RETURNING creation_date, last_modified 
	`
	var row pgx.Row
	if len(tx) == 0 {
		row = DB.QueryRow(
			context.Background(),
			insertQuery,
			newProject.ID,
			newProject.Name,
			newProject.Owner,
		)
	} else {
		row = tx[0].QueryRow(
			context.Background(),
			insertQuery,
			newProject.ID,
			newProject.Name,
			newProject.Owner,
		)
	}

	err := row.Scan(
		&newProject.CreationDate,
		&newProject.LastModified,
	)
	if err != nil {
		fmt.Printf("error while inserting project into db %+v \n", err)
		return newProject, err
	}
	return newProject, nil
}

func DeleteProject(projectID string, email string, tx ...pgx.Tx) error {
	deleteQuery := `
	DELETE FROM projects
	WHERE id = $1 AND owner = $2
	`
	var result pgconn.CommandTag
	var err error
	if len(tx) == 0 {
		result, err = DB.Exec(
			context.Background(),
			deleteQuery,
			projectID,
			email,
		)
	} else {
		result, err = tx[0].Exec(
			context.Background(),
			deleteQuery,
			projectID,
			email,
		)
	}
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

func GetTreeIDsForProject(projID string, tx ...pgx.Tx) ([]string, error) {
	selectQuery := `
	SELECT id FROM trees
	WHERE project_id = $1
	`
	var treeIDs []string
	var rows pgx.Rows
	var err error

	if len(tx) == 0 {
		rows, err = DB.Query(context.Background(), selectQuery, projID)
	} else {
		rows, err = tx[0].Query(context.Background(), selectQuery, projID)
	}
	if err != nil {
		return treeIDs, err
	}
	defer rows.Close()

	for rows.Next() {
		var treeID string
		err := rows.Scan(
			&treeID,
		)
		if err != nil {
			fmt.Printf("error scanning treeIDs: %+v\n", err)
			return nil, err
		}
		treeIDs = append(treeIDs, treeID)
	}

	return treeIDs, nil
}

func MatchProjectWithEmail(projID string, email string, tx ...pgx.Tx) (bool, error) {
	selectQuery := `
	SELECT EXISTS (
		SELECT 1
		FROM projects
		WHERE id = $1 AND owner = $2
	);
	`
	var found bool
	var err error

	if len(tx) == 0 {
		err = DB.QueryRow(context.Background(), selectQuery, projID, email).Scan(&found)
	} else {
		err = tx[0].QueryRow(context.Background(), selectQuery, projID, email).Scan(&found)
	}
	if err != nil {
		return false, err
	}
	return found, err
}
