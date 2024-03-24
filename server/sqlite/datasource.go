package sqlite

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

	txManager beans.TxManager
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

func (ds *datasource) TxManager() beans.TxManager {
	return ds.txManager
}

func NewDataSource(pool *Pool) *datasource {
	return &datasource{
		accountRepository:       &accountRepository{repository{pool}},
		budgetRepository:        &budgetRepository{repository{pool}},
		categoryRepository:      &categoryRepository{repository{pool}},
		monthRepository:         &monthRepository{repository{pool}},
		monthCategoryRepository: &monthCategoryRepository{repository{pool}},
		payeeRepository:         &payeeRepository{repository{pool}},
		transactionRepository:   &TransactionRepository{repository{pool}},
		userRepository:          &userRepository{repository{pool}},

		txManager: &txManager{pool},
	}
}
