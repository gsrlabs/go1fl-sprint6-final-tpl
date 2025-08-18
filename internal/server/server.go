package server

import (
	"log"
	"net/http"

	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
	"github.com/go-chi/chi"
)

type Server struct {
	httpServer *http.Server
	logger     *log.Logger // думаю лишнее
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(logger *log.Logger) *Server {
	r := chi.NewRouter()
	r.Get("/", handlers.HandleMain)
	r.Post("/upload", handlers.HandleUpload)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":8080",
			Handler:      r,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening on the configured server address.
func (s *Server) Start() error {
	s.logger.Println("Starting server on :8080")
	return s.httpServer.ListenAndServe()
}
