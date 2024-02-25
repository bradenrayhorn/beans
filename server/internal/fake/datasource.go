package fake

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

func NewDataSource() *datasource {
	database := &database{
		accounts:        make(map[beans.ID]beans.Account),
		budgets:         make(map[beans.ID]beans.Budget),
		budgetUsers:     make(map[beans.ID][]beans.ID),
		categories:      make(map[beans.ID]beans.Category),
		categoryGroups:  make(map[beans.ID]beans.CategoryGroup),
		months:          make(map[beans.ID]beans.Month),
		monthCategories: make(map[beans.ID]beans.MonthCategory),
		payees:          make(map[beans.ID]beans.Payee),
		transactions:    make(map[beans.ID]beans.Transaction),
		users:           make(map[beans.ID]beans.User),
	}
	return &datasource{
		accountRepository:       &accountRepository{repository{database}},
		budgetRepository:        &budgetRepository{repository{database}},
		categoryRepository:      &categoryRepository{repository{database}},
		monthRepository:         &monthRepository{repository{database}},
		monthCategoryRepository: &monthCategoryRepository{repository{database}},
		payeeRepository:         &payeeRepository{repository{database}},
		transactionRepository:   &transactionRepository{repository{database}},
		userRepository:          &userRepository{repository{database}},

		txManager: &txManager{},
	}
}
