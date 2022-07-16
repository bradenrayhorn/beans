package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	sv     *http.Server

	userService beans.UserService
}

func NewServer(us beans.UserService) *Server {
	s := &Server{
		router:      chi.NewRouter(),
		sv:          &http.Server{},
		userService: us,
	}

	s.sv.Handler = s.router

	s.router.Get("/health-check", s.handleHealthCheck)
	s.router.Route("/api/v1", func(r chi.Router) {

		r.Route("/user", func(r chi.Router) {
			r.Post("/register", s.handleUserRegister())
			r.Post("/login", s.handleUserLogin())
		})

	})

	return s
}

func (s *Server) Open() error {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		return err
	}

	go s.sv.Serve(ln)

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return s.sv.Shutdown(ctx)
}
