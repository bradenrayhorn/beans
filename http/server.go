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
	router    *chi.Mux
	sv        *http.Server
	boundAddr string

	accountRepository       beans.AccountRepository
	accountService          beans.AccountService
	budgetRepository        beans.BudgetRepository
	budgetService           beans.BudgetService
	categoryRepository      beans.CategoryRepository
	categoryService         beans.CategoryService
	monthCategoryRepository beans.MonthCategoryRepository
	monthRepository         beans.MonthRepository
	monthService            beans.MonthService
	sessionRepository       beans.SessionRepository
	transactionRepository   beans.TransactionRepository
	transactionService      beans.TransactionService
	userRepository          beans.UserRepository
	userService             beans.UserService
}

func NewServer(
	ar beans.AccountRepository,
	as beans.AccountService,
	br beans.BudgetRepository,
	bs beans.BudgetService,
	cr beans.CategoryRepository,
	cs beans.CategoryService,
	mcr beans.MonthCategoryRepository,
	mr beans.MonthRepository,
	ms beans.MonthService,
	sr beans.SessionRepository,
	tr beans.TransactionRepository,
	ts beans.TransactionService,
	ur beans.UserRepository,
	us beans.UserService,
) *Server {
	s := &Server{
		router:                  chi.NewRouter(),
		sv:                      &http.Server{},
		accountRepository:       ar,
		accountService:          as,
		budgetRepository:        br,
		budgetService:           bs,
		categoryRepository:      cr,
		categoryService:         cs,
		monthCategoryRepository: mcr,
		monthRepository:         mr,
		monthService:            ms,
		sessionRepository:       sr,
		transactionRepository:   tr,
		transactionService:      ts,
		userRepository:          ur,
		userService:             us,
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
				r.Post("/", s.handleAccountCreate())
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", s.handleTransactionGetAll())
				r.Post("/", s.handleTransactionCreate())
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", s.handleCategoryGetAll())
				r.Post("/", s.handleCategoryCreate())

				r.Route("/groups", func(r chi.Router) {
					r.Post("/", s.handleCategoryGroupCreate())
				})
			})

			r.Route("/months", func(r chi.Router) {
				r.Get("/{id}", s.handleMonthGet())
			})

		})
	})

	s.router.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 api"))
	})

	s.router.Get("/*", s.handleServeFrontend())

	return s
}

func (s *Server) Open(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.boundAddr = ln.Addr().String()

	go s.sv.Serve(ln)

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
