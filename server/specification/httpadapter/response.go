package httpadapter

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

func mapAll[T any, K any](objs []T, mapper func(T) K) []K {
	var models []K
	for _, m := range objs {
		mapped := mapper(m)

		models = append(models, mapped)
	}

	return models
}

func mapBudget(t response.Budget) beans.Budget {
	return beans.Budget{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

func mapCategory(t response.Category) beans.Category {
	return beans.Category{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

func mapCategoryGroup(t response.CategoryGroup) beans.CategoryGroup {
	return beans.CategoryGroup{
		ID:       t.ID,
		Name:     beans.Name(t.Name),
		IsIncome: t.IsIncome,
	}
}

func mapListAccount(t response.ListAccount) beans.Account {
	return beans.Account{
		ID:      t.ID,
		Name:    beans.Name(t.Name),
		Balance: t.Balance,
	}
}

func mapTransaction(t response.Transaction) beans.Transaction {
	transaction := beans.Transaction{
		ID:        t.ID,
		AccountID: t.Account.ID,
		Account: &beans.Account{
			ID:   t.Account.ID,
			Name: t.Account.Name,
		},
		Amount: t.Amount,
		Date:   t.Date,
		Notes:  t.Notes,
	}

	if t.Category != nil {
		transaction.CategoryID = t.Category.ID
		transaction.CategoryName = beans.NewNullString(string(t.Category.Name))
	}
	if t.Payee != nil {
		transaction.PayeeID = t.Payee.ID
		transaction.PayeeName = beans.NewNullString(string(t.Payee.Name))
	}

	return transaction
}
