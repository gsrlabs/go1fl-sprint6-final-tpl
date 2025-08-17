package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stderr, "SERVER: ", log.LstdFlags|log.Lshortfile)
	srv := server.NewServer(logger)

	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
