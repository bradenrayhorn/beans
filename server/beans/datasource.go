package beans

// A collection of repositories that represents the primary datastore of beans.
type DataSource interface {
	AccountRepository() AccountRepository
	BudgetRepository() BudgetRepository
	CategoryRepository() CategoryRepository
	MonthRepository() MonthRepository
	MonthCategoryRepository() MonthCategoryRepository
	PayeeRepository() PayeeRepository
	TransactionRepository() TransactionRepository
	UserRepository() UserRepository

	TxManager() TxManager
}
