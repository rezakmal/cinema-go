package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectDB() (*sql.DB, error) {
	// get database connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// set default values if not provided
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "cinema_db"
	}

	// build connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// verify connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Connected to database")
	return db, nil
}

// RunMigrations executes database migration scripts
func RunMigrations(db *sql.DB) error {
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS cinema (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			location TEXT NOT NULL,
			rating DECIMAL(3,1) DEFAULT 0.0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_cinema_name ON cinema(name);
		CREATE INDEX IF NOT EXISTS idx_cinema_location ON cinema(location);
	`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
