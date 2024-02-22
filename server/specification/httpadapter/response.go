package httpadapter

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

func mapAll[T any, K any](objs []T, mapper func(T) K) []K {
	models := []K{}
	for _, m := range objs {
		mapped := mapper(m)

		models = append(models, mapped)
	}

	return models
}

// account

func mapAccount(t response.Account) beans.Account {
	return beans.Account{ID: t.ID, Name: beans.Name(t.Name)}
}

func mapListAccount(t response.ListAccount) beans.AccountWithBalance {
	return beans.AccountWithBalance{
		Account: beans.Account{ID: t.ID, Name: beans.Name(t.Name)},
		Balance: t.Balance,
	}
}

// budget

func mapBudget(t response.Budget) beans.Budget {
	return beans.Budget{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

// category

func mapCategory(t response.Category) beans.Category {
	return beans.Category{
		ID:      t.ID,
		Name:    beans.Name(t.Name),
		GroupID: t.GroupID,
	}
}

func mapRelatedCategory(t response.AssociatedCategory) beans.RelatedCategory {
	return beans.RelatedCategory{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

func mapCategoryGroupWithCategories(t response.CategoryGroup) beans.CategoryGroupWithCategories {
	return beans.CategoryGroupWithCategories{
		CategoryGroup: beans.CategoryGroup{
			ID:       t.ID,
			Name:     beans.Name(t.Name),
			IsIncome: t.IsIncome,
		},
		Categories: mapAll(t.Categories, mapRelatedCategory),
	}
}

// month

func mapMonthCategory(t response.MonthCategory) beans.MonthCategoryWithDetails {
	return beans.MonthCategoryWithDetails{
		ID:         t.ID,
		CategoryID: t.CategoryID,
		Amount:     t.Assigned,
		Activity:   t.Activity,
		Available:  t.Available,
	}
}

func mapMonthWithDetails(t response.Month) beans.MonthWithDetails {
	return beans.MonthWithDetails{
		Month: beans.Month{
			ID:        t.ID,
			Date:      t.Date,
			Carryover: t.Carryover,
		},
		CarriedOver: t.CarriedOver,
		Income:      t.Income,
		Assigned:    t.Assigned,
		Budgetable:  t.Budgetable,
		Categories:  mapAll(t.Categories, mapMonthCategory),
	}
}

// payee

func mapPayee(t response.Payee) beans.Payee {
	return beans.Payee{ID: t.ID, Name: beans.Name(t.Name)}
}

// transaction

func mapTransactionWithRelations(t response.Transaction) beans.TransactionWithRelations {
	transaction := beans.TransactionWithRelations{
		Transaction: beans.Transaction{
			ID:        t.ID,
			AccountID: t.Account.ID,
			Amount:    t.Amount,
			Date:      t.Date,
			Notes:     t.Notes,
		},
		Account: beans.RelatedAccount{
			ID:   t.Account.ID,
			Name: t.Account.Name,
		},
	}

	if t.Category != nil {
		transaction.Category = beans.OptionalWrap(beans.RelatedCategory{ID: t.Category.ID, Name: t.Category.Name})
	}
	if t.Payee != nil {
		transaction.Payee = beans.OptionalWrap(beans.RelatedPayee{ID: t.Payee.ID, Name: t.Payee.Name})
	}

	return transaction
}
