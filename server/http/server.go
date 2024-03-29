package http

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/service"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router    *chi.Mux
	sv        *http.Server
	boundAddr string

	contracts *contract.Contracts
	services  *service.All
}

func NewServer(
	contracts *contract.Contracts,
	services *service.All,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		sv:     &http.Server{},

		contracts: contracts,
		services:  services,
	}

	s.sv.Handler = s.router

	s.router.Get("/health-check", s.handleHealthCheck)
	s.router.Route("/api/v1", func(r chi.Router) {

		r.Route("/user", func(r chi.Router) {
			r.Post("/register", s.handleUserRegister())
			r.Post("/login", s.handleUserLogin())
			r.Group(func(r chi.Router) {
				r.Use(s.authenticate)
				r.Get("/me", s.handleUserMe())
				r.Post("/logout", s.handleUserLogout())
			})
		})

		r.Route("/budgets", func(r chi.Router) {
			r.Use(s.authenticate)
			r.Post("/", s.handleBudgetCreate())
			r.Get("/", s.handleBudgetGetAll())
			r.Get("/{budgetID}", s.handleBudgetGet())
		})

		// endpoints that require budget header
		r.Group(func(r chi.Router) {
			r.Use(s.authenticate, s.parseBudgetHeader)

			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", s.handleAccountsGet())
				r.Get("/transactable", s.handleAccountsGetTransactable())

				r.Post("/", s.handleAccountCreate())
				r.Get("/{accountID}", s.handleAccountGet())
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", s.handleCategoryGetAll())
				r.Post("/", s.handleCategoryCreate())
				r.Get("/{categoryID}", s.handleCategoryGetCategory())

				r.Route("/groups", func(r chi.Router) {
					r.Post("/", s.handleCategoryGroupCreate())
					r.Get("/{categoryGroupID}", s.handleCategoryGetCategoryGroup())
				})
			})

			r.Route("/months", func(r chi.Router) {
				r.Route("/{monthID}", func(r chi.Router) {
					r.Put("/", s.handleMonthUpdate())
					r.Post("/categories", s.handleMonthCategoryUpdate())
				})

				r.Get("/{date}", s.handleMonthGetOrCreate())
			})

			r.Route("/payees", func(r chi.Router) {
				r.Get("/", s.handlePayeeGetAll())
				r.Post("/", s.handlePayeeCreate())
				r.Get("/{payeeID}", s.handlePayeeGet())
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", s.handleTransactionGetAll())
				r.Post("/", s.handleTransactionCreate())
				r.Post("/delete", s.handleTransactionDelete())
				r.Put("/{transactionID}", s.handleTransactionUpdate())
				r.Get("/{transactionID}", s.handleTransactionGet())
				r.Get("/{transactionID}/splits", s.handleTransactionGetSplits())
			})

		})
	})

	s.router.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return s
}

func (s *Server) Open(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.boundAddr = ln.Addr().String()

	go func() {
		err := s.sv.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
		}
	}()

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return s.sv.Shutdown(ctx)
}

func (s *Server) GetBoundAddr() string {
	return s.boundAddr
}
