package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	httpServer *http.Server
	pool       *pgxpool.Pool

	config Config

	txManager beans.TxManager

	accountRepository       beans.AccountRepository
	budgetRepository        beans.BudgetRepository
	categoryRepository      beans.CategoryRepository
	monthRepository         beans.MonthRepository
	monthCategoryRepository beans.MonthCategoryRepository
	payeeRepository         beans.PayeeRepository
	sessionRepository       beans.SessionRepository
	transactionRepository   beans.TransactionRepository
	userRepository          beans.UserRepository
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

	a.txManager = postgres.NewTxManager(pool)

	a.sessionRepository = inmem.NewSessionRepository()

	a.accountRepository = postgres.NewAccountRepository(pool)
	a.budgetRepository = postgres.NewBudgetRepository(pool)
	a.categoryRepository = postgres.NewCategoryRepository(pool)
	a.monthRepository = postgres.NewMonthRepository(pool)
	a.monthCategoryRepository = postgres.NewMonthCategoryRepository(pool)
	a.payeeRepository = postgres.NewPayeeRepository(pool)
	a.transactionRepository = postgres.NewTransactionRepository(pool)
	a.userRepository = postgres.NewUserRepository(pool)

	a.httpServer = http.NewServer(
		a.accountRepository,
		a.budgetRepository,
		a.categoryRepository,
		a.monthRepository,
		a.monthCategoryRepository,
		a.sessionRepository,
		a.transactionRepository,
		a.userRepository,

		contract.NewAccountContract(a.accountRepository),
		contract.NewBudgetContract(a.budgetRepository, a.categoryRepository, a.monthRepository, a.txManager),
		contract.NewCategoryContract(
			a.categoryRepository,
			a.monthCategoryRepository,
			a.monthRepository,
			a.txManager,
		),
		contract.NewMonthContract(
			a.categoryRepository,
			a.monthRepository,
			a.monthCategoryRepository,
			a.transactionRepository,
			a.txManager,
		),
		contract.NewPayeeContract(a.payeeRepository),
		contract.NewTransactionContract(
			a.transactionRepository,
			a.accountRepository,
			a.categoryRepository,
			a.monthCategoryRepository,
			a.monthRepository,
			a.payeeRepository,
		),
		contract.NewUserContract(
			a.sessionRepository,
			a.userRepository,
		),
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

func (a *Application) TransactionRepository() beans.TransactionRepository {
	return a.transactionRepository
}

func (a *Application) TxManager() beans.TxManager {
	return a.txManager
}