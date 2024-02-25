package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type budgetRepository struct{ repository }

var _ beans.BudgetRepository = (*budgetRepository)(nil)

func (r *budgetRepository) Create(ctx context.Context, tx beans.Tx, id beans.ID, name beans.Name, userID beans.ID) error {
	r.acquire(func() { r.database.budgetsMU.RLock() })
	if _, ok := r.database.budgets[id]; ok {
		r.database.budgetsMU.RUnlock()
		return errors.New("duplicate")
	} else {
		r.database.budgetsMU.RUnlock()
	}

	r.txOrNow(tx, func() {
		r.acquire(func() { r.database.budgetsMU.Lock() })
		defer r.database.budgetsMU.Unlock()

		r.database.budgets[id] = beans.Budget{ID: id, Name: name}
		r.database.budgetUsers[id] = []beans.ID{userID}
	})

	return nil
}

func (r *budgetRepository) Get(ctx context.Context, id beans.ID) (beans.Budget, error) {
	r.acquire(func() { r.database.budgetsMU.RLock() })
	defer r.database.budgetsMU.RUnlock()

	if v, ok := r.database.budgets[id]; ok {
		return v, nil
	}

	return beans.Budget{}, beans.NewError(beans.ENOTFOUND, "budget not found")
}

func (r *budgetRepository) GetBudgetsForUser(ctx context.Context, userID beans.ID) ([]beans.Budget, error) {
	r.acquire(func() { r.database.budgetsMU.RLock() })
	defer r.database.budgetsMU.RUnlock()

	budgets := make([]beans.Budget, 0)
	for k, v := range r.database.budgetUsers {
		if find(v, func(it beans.ID) bool { return it == userID }) != nil {
			budgets = append(budgets, r.database.budgets[k])
		}
	}

	return budgets, nil
}

func (r *budgetRepository) GetBudgetUserIDs(ctx context.Context, id beans.ID) ([]beans.ID, error) {
	r.acquire(func() { r.database.budgetsMU.RLock() })
	defer r.database.budgetsMU.RUnlock()

	if v, ok := r.database.budgetUsers[id]; ok {
		return v, nil
	}

	return make([]beans.ID, 0), nil
}
