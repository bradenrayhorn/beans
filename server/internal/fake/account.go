package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type accountRepository struct{ repository }

var _ beans.AccountRepository = (*accountRepository)(nil)

func (r *accountRepository) Create(ctx context.Context, id beans.ID, name beans.Name, budgetID beans.ID) error {
	r.acquire(func() { r.database.accountsMU.Lock() })
	defer r.database.accountsMU.Unlock()

	if _, ok := r.database.accounts[id]; ok {
		return errors.New("duplicate")
	}

	r.database.accounts[id] = beans.Account{
		ID:       id,
		Name:     name,
		BudgetID: budgetID,
	}

	return nil
}

func (r *accountRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Account, error) {
	r.acquire(func() { r.database.accountsMU.RLock() })
	defer r.database.accountsMU.RUnlock()

	if account, ok := r.database.accounts[id]; ok {
		if account.BudgetID == budgetID {
			return account, nil
		}
	}

	return beans.Account{}, beans.NewError(beans.ENOTFOUND, "account not found")
}

func (r *accountRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.AccountWithBalance, error) {
	r.acquire(func() {
		r.database.accountsMU.RLock()
		r.database.transactionsMU.RLock()
	})
	defer r.database.accountsMU.RUnlock()
	defer r.database.transactionsMU.RUnlock()

	accounts := filter(
		values(r.database.accounts), func(it beans.Account) bool {
			return it.BudgetID == budgetID
		})

	return mapVals(accounts, func(it beans.Account) beans.AccountWithBalance {
		transactions := filter(values(r.database.transactions), func(t beans.Transaction) bool {
			return t.AccountID == it.ID
		})

		balance := reduce(transactions, beans.NewAmount(0, 0), func(it beans.Transaction, acc beans.Amount) beans.Amount {
			r, err := beans.Arithmetic.Add(acc, it.Amount)
			if err != nil {
				panic(err)
			}
			return r
		})

		return beans.AccountWithBalance{
			Account: it,
			Balance: balance,
		}
	}), nil
}
