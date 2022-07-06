package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.router.Get("/health-check", s.handleHealthCheck)

	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":8000", s.router)
}
