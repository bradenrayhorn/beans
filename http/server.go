package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router    *chi.Mux
	sv        *http.Server
	boundAddr string

	accountContract     beans.AccountContract
	budgetContract      beans.BudgetContract
	categoryContract    beans.CategoryContract
	monthContract       beans.MonthContract
	payeeContract       beans.PayeeContract
	transactionContract beans.TransactionContract
	userContract        beans.UserContract

	accountRepository       beans.AccountRepository
	budgetRepository        beans.BudgetRepository
	categoryRepository      beans.CategoryRepository
	monthRepository         beans.MonthRepository
	monthCategoryRepository beans.MonthCategoryRepository
	sessionRepository       beans.SessionRepository
	transactionRepository   beans.TransactionRepository
	userRepository          beans.UserRepository
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

	accountContract beans.AccountContract,
	budgetContract beans.BudgetContract,
	categoryContract beans.CategoryContract,
	monthContract beans.MonthContract,
	payeeContract beans.PayeeContract,
	transactionContract beans.TransactionContract,
	userContract beans.UserContract,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		sv:     &http.Server{},

		accountContract:     accountContract,
		budgetContract:      budgetContract,
		categoryContract:    categoryContract,
		monthContract:       monthContract,
		payeeContract:       payeeContract,
		transactionContract: transactionContract,
		userContract:        userContract,

		accountRepository:       ar,
		budgetRepository:        br,
		categoryRepository:      cr,
		monthRepository:         mr,
		monthCategoryRepository: mcr,
		sessionRepository:       sr,
		transactionRepository:   tr,
		userRepository:          ur,
	}

	s.sv.Handler = s.router
	s.router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "OPTIONS"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowCredentials: true,

		AllowedHeaders: []string{"Budget-ID", "Content-Type"},
	}))

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

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", s.handleCategoryGetAll())
				r.Post("/", s.handleCategoryCreate())

				r.Route("/groups", func(r chi.Router) {
					r.Post("/", s.handleCategoryGroupCreate())
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
				r.Get("/", s.handlePayeeCreate())
				r.Post("/", s.handlePayeeGetAll())
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", s.handleTransactionGetAll())
				r.Post("/", s.handleTransactionCreate())
				r.Put("/{transactionID}", s.handleTransactionUpdate())
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
			log.Err(err).Msg("http server error")
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
