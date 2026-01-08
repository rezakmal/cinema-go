package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectDB(*sql.DB) {
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
