package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/http"
	"github.com/bradenrayhorn/beans/inmem"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	httpServer *http.Server
	pool       *pgxpool.Pool

	config Config

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

func NewApplication(c Config) *Application {
	return &Application{
		config: c,
	}
}

func (a *Application) Start() error {
	pool, err := postgres.CreatePool(
		fmt.Sprintf("postgres://%s:%s@%s/%s",
			a.config.Postgres.Username,
			a.config.Postgres.Password,
			a.config.Postgres.Addr,
			a.config.Postgres.Database,
		))

	if err != nil {
		panic(err)
	}
	a.pool = pool

	a.accountRepository = postgres.NewAccountRepository(pool)
	a.accountService = logic.NewAccountService(a.accountRepository)
	a.budgetRepository = postgres.NewBudgetRepository(pool)
	a.budgetService = logic.NewBudgetService(a.budgetRepository)
	a.categoryRepository = postgres.NewCategoryRepository(pool)
	a.categoryService = logic.NewCategoryService(a.categoryRepository)
	a.monthCategoryRepository = postgres.NewMonthCategoryRepository(pool)
	a.monthRepository = postgres.NewMonthRepository(pool)
	a.monthService = logic.NewMonthService(a.monthRepository)
	a.sessionRepository = inmem.NewSessionRepository()
	a.transactionRepository = postgres.NewTransactionRepository(pool)
	a.transactionService = logic.NewTransactionService(a.transactionRepository, a.accountRepository)
	a.userRepository = postgres.NewUserRepository(pool)
	a.userService = &logic.UserService{UserRepository: a.userRepository}

	a.httpServer = http.NewServer(
		a.accountRepository,
		a.accountService,
		a.budgetRepository,
		a.budgetService,
		a.categoryRepository,
		a.categoryService,
		a.monthCategoryRepository,
		a.monthRepository,
		a.monthService,
		a.sessionRepository,
		a.transactionRepository,
		a.transactionService,
		a.userRepository,
		a.userService,
	)
	if err := a.httpServer.Open(":" + a.config.Port); err != nil {
		panic(err)
	}

	return nil
}

func (a *Application) Stop() error {
	if err := a.httpServer.Close(); err != nil {
		return err
	}

	a.pool.Close()

	return nil
}

func (a *Application) HttpServer() *http.Server {
	return a.httpServer
}

func (a *Application) AccountRepository() beans.AccountRepository {
	return a.accountRepository
}

func (a *Application) BudgetRepository() beans.BudgetRepository {
	return a.budgetRepository
}

func (a *Application) CategoryRepository() beans.CategoryRepository {
	return a.categoryRepository
}

func (a *Application) MonthRepository() beans.MonthRepository {
	return a.monthRepository
}

func (a *Application) MonthCategoryRepository() beans.MonthCategoryRepository {
	return a.monthCategoryRepository
}

func (a *Application) UserRepository() beans.UserRepository {
	return a.userRepository
}

func (a *Application) SessionRepository() beans.SessionRepository {
	return a.sessionRepository
}
