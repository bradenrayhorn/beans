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

	accountContract     beans.AccountContract
	budgetContract      beans.BudgetContract
	categoryContract    beans.CategoryContract
	monthContract       beans.MonthContract
	transactionContract beans.TransactionContract

	accountRepository       beans.AccountRepository
	budgetRepository        beans.BudgetRepository
	categoryRepository      beans.CategoryRepository
	monthRepository         beans.MonthRepository
	monthCategoryRepository beans.MonthCategoryRepository
	sessionRepository       beans.SessionRepository
	transactionRepository   beans.TransactionRepository
	userRepository          beans.UserRepository
	userService             beans.UserService
}

func NewServer(
	ar beans.AccountRepository,
	br beans.BudgetRepository,
	cr beans.CategoryRepository,
	mr beans.MonthRepository,
	mcr beans.MonthCategoryRepository,
	sr beans.SessionRepository,
	tr beans.TransactionRepository,
	ur beans.UserRepository,
	us beans.UserService,

	accountContract beans.AccountContract,
	budgetContract beans.BudgetContract,
	categoryContract beans.CategoryContract,
	monthContract beans.MonthContract,
	transactionContract beans.TransactionContract,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		sv:     &http.Server{},

		accountContract:     accountContract,
		budgetContract:      budgetContract,
		categoryContract:    categoryContract,
		monthContract:       monthContract,
		transactionContract: transactionContract,

		accountRepository:       ar,
		budgetRepository:        br,
		categoryRepository:      cr,
		monthRepository:         mr,
		monthCategoryRepository: mcr,
		sessionRepository:       sr,
		transactionRepository:   tr,
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
				r.Post("/", s.handleMonthCreate())

				r.Route("/{monthID}", func(r chi.Router) {
					r.Use(s.validateMonth)

					r.Get("/", s.handleMonthGet())
					r.Post("/categories", s.handleMonthCategoryUpdate())
				})
			})

		})
	})

	s.router.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
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
