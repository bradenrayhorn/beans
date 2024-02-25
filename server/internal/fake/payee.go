package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type payeeRepository struct{ repository }

var _ beans.PayeeRepository = (*payeeRepository)(nil)

func (r *payeeRepository) Create(ctx context.Context, payee beans.Payee) error {
	r.acquire(func() { r.database.payeesMU.Lock() })
	defer r.database.payeesMU.Unlock()

	if _, ok := r.database.payees[payee.ID]; ok {
		return errors.New("duplicate")
	}

	r.database.payees[payee.ID] = payee

	return nil
}

func (r *payeeRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Payee, error) {
	r.acquire(func() { r.database.payeesMU.RLock() })
	defer r.database.payeesMU.RUnlock()

	if payee, ok := r.database.payees[id]; ok {
		if payee.BudgetID == budgetID {
			return payee, nil
		}
	}

	return beans.Payee{}, beans.NewError(beans.ENOTFOUND, "payee not found")
}

func (r *payeeRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Payee, error) {
	r.acquire(func() { r.database.payeesMU.RLock() })
	defer r.database.payeesMU.RUnlock()

	return filter(values(r.database.payees), func(it beans.Payee) bool { return it.BudgetID == budgetID }), nil
}
