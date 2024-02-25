package fake

import (
	"sync"

	"github.com/bradenrayhorn/beans/server/beans"
)

type database struct {
	accounts   map[beans.ID]beans.Account
	accountsMU sync.RWMutex

	budgets     map[beans.ID]beans.Budget
	budgetUsers map[beans.ID][]beans.ID
	budgetsMU   sync.RWMutex

	categories     map[beans.ID]beans.Category
	categoryGroups map[beans.ID]beans.CategoryGroup
	categoriesMU   sync.RWMutex

	months   map[beans.ID]beans.Month
	monthsMU sync.RWMutex

	monthCategories   map[beans.ID]beans.MonthCategory
	monthCategoriesMU sync.RWMutex

	payees   map[beans.ID]beans.Payee
	payeesMU sync.RWMutex

	transactions   map[beans.ID]beans.Transaction
	transactionsMU sync.RWMutex

	users   map[beans.ID]beans.User
	usersMU sync.RWMutex

	mu sync.RWMutex
}
