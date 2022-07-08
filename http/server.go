package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux

	userService beans.UserService
}

func NewServer(us beans.UserService) *Server {
	s := &Server{
		router:      chi.NewRouter(),
		userService: us,
	}

	s.router.Get("/health-check", s.handleHealthCheck)
	s.router.Route("/api/v1", func(r chi.Router) {

		r.Route("/user", func(r chi.Router) {
			r.Post("/register", s.handleUserRegister())
			r.Post("/login", s.handleUserLogin())
		})

	})

	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":8000", s.router)
}
