package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/abhinavmsra/go-api/internal/app"
	"github.com/abhinavmsra/go-api/internal/repository"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	connectionString := "postgres://api@db:5432/api_development?sslmode=disable"

	// setup database
	db, err := setupDatabase(connectionString)
	if err != nil {
		return err
	}

	storage := repository.NewStorage(db)

	err = storage.RunMigrations(connectionString)

	if err != nil {
		panic(err)
	}

	return nil

	// setup app server
	server := app.NewServer()
	return server.Run()
}

func setupDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}