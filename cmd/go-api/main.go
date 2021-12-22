package main

import (
	"log"
	"os"

	"github.com/abhinavmsra/go-api/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	server := app.NewServer()
	return server.Run()
}
