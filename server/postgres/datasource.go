package postgres

import "github.com/bradenrayhorn/beans/server/beans"

var _ beans.DataSource = (*datasource)(nil)

type datasource struct {
	accountRepository       beans.AccountRepository
	budgetRepository        beans.BudgetRepository
	categoryRepository      beans.CategoryRepository
	monthRepository         beans.MonthRepository
	monthCategoryRepository beans.MonthCategoryRepository
	payeeRepository         beans.PayeeRepository
	transactionRepository   beans.TransactionRepository
	userRepository          beans.UserRepository
}

func (ds *datasource) AccountRepository() beans.AccountRepository {
	return ds.accountRepository
}

func (ds *datasource) BudgetRepository() beans.BudgetRepository {
	return ds.budgetRepository
}

func (ds *datasource) CategoryRepository() beans.CategoryRepository {
	return ds.categoryRepository
}

func (ds *datasource) MonthRepository() beans.MonthRepository {
	return ds.monthRepository
}

func (ds *datasource) MonthCategoryRepository() beans.MonthCategoryRepository {
	return ds.monthCategoryRepository
}

func (ds *datasource) PayeeRepository() beans.PayeeRepository {
	return ds.payeeRepository
}

func (ds *datasource) TransactionRepository() beans.TransactionRepository {
	return ds.transactionRepository
}

func (ds *datasource) UserRepository() beans.UserRepository {
	return ds.userRepository
}

func NewDataSource(pool *DbPool) *datasource {
	return &datasource{
		accountRepository:       NewAccountRepository(pool),
		budgetRepository:        NewBudgetRepository(pool),
		categoryRepository:      NewCategoryRepository(pool),
		monthRepository:         NewMonthRepository(pool),
		monthCategoryRepository: NewMonthCategoryRepository(pool),
		payeeRepository:         NewPayeeRepository(pool),
		transactionRepository:   NewTransactionRepository(pool),
		userRepository:          NewUserRepository(pool),
	}
}
