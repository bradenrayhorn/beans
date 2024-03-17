package fake

import (
	"context"
	"errors"
	"slices"

	"github.com/bradenrayhorn/beans/server/beans"
)

type transactionRepository struct{ repository }

var _ beans.TransactionRepository = (*transactionRepository)(nil)

func (r *transactionRepository) Create(ctx context.Context, transactions []beans.Transaction) error {
	r.acquire(func() { r.database.transactionsMU.Lock() })
	defer r.database.transactionsMU.Unlock()

	for _, transaction := range transactions {
		if _, ok := r.database.transactions[transaction.ID]; ok {
			return errors.New("duplicate")
		}

		r.database.transactions[transaction.ID] = transaction
	}

	return nil
}

func (r *transactionRepository) Update(ctx context.Context, transactions []beans.Transaction) error {
	r.acquire(func() { r.database.transactionsMU.Lock() })
	defer r.database.transactionsMU.Unlock()

	for _, transaction := range transactions {
		if _, ok := r.database.transactions[transaction.ID]; ok {
			r.database.transactions[transaction.ID] = transaction
		}
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, budgetID beans.ID, transactionIDs []beans.ID) error {
	r.acquire(func() {
		r.database.transactionsMU.Lock()
		r.database.accountsMU.RLock()
	})
	defer r.database.transactionsMU.Unlock()
	defer r.database.accountsMU.RUnlock()

	transactions := filter(values(r.database.transactions), func(it beans.Transaction) bool {
		return slices.Contains(transactionIDs, it.ID)
	})

	for _, v := range transactions {
		if account, ok := r.database.accounts[v.AccountID]; ok {
			if account.BudgetID == budgetID {
				delete(r.database.transactions, v.ID)

				// emulate foreign key on cascade delete
				if !v.TransferID.Empty() {
					delete(r.database.transactions, v.TransferID)
				}
			}
		}
	}

	return nil
}

func (r *transactionRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Transaction, error) {
	r.acquire(func() {
		r.database.transactionsMU.RLock()
		r.database.accountsMU.RLock()
	})
	defer r.database.transactionsMU.RUnlock()
	defer r.database.accountsMU.RUnlock()

	if t, ok := r.database.transactions[id]; ok {

		if account, ok := r.database.accounts[t.AccountID]; ok {
			if account.BudgetID == budgetID {
				return t, nil
			}
		}
	}

	return beans.Transaction{}, beans.NewError(beans.ENOTFOUND, "not found")
}

func (r *transactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.TransactionWithRelations, error) {
	r.acquire(func() {
		r.database.transactionsMU.RLock()
		r.database.accountsMU.RLock()
		r.database.categoriesMU.RLock()
		r.database.payeesMU.RLock()
	})
	defer r.database.transactionsMU.RUnlock()
	defer r.database.accountsMU.RUnlock()
	defer r.database.categoriesMU.RUnlock()
	defer r.database.payeesMU.RUnlock()

	transactions := filter(values(r.database.transactions), func(it beans.Transaction) bool {
		if account, ok := r.database.accounts[it.AccountID]; ok {
			if account.BudgetID == budgetID {
				return true
			}
		}
		return false
	})

	return mapVals(transactions, r.mapWithRelations), nil
}

func (r *transactionRepository) GetWithRelations(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.TransactionWithRelations, error) {
	r.acquire(func() {
		r.database.transactionsMU.RLock()
		r.database.accountsMU.RLock()
		r.database.categoriesMU.RLock()
		r.database.payeesMU.RLock()
	})
	defer r.database.transactionsMU.RUnlock()
	defer r.database.accountsMU.RUnlock()
	defer r.database.categoriesMU.RUnlock()
	defer r.database.payeesMU.RUnlock()

	if t, ok := r.database.transactions[id]; ok {

		if account, ok := r.database.accounts[t.AccountID]; ok {
			if account.BudgetID == budgetID {
				return r.mapWithRelations(t), nil
			}
		}
	}

	return beans.TransactionWithRelations{}, beans.NewError(beans.ENOTFOUND, "not found")
}

func (r *transactionRepository) GetActivityByCategory(ctx context.Context, budgetID beans.ID, from beans.Date, to beans.Date) (map[beans.ID]beans.Amount, error) {
	r.acquire(func() {
		r.database.transactionsMU.RLock()
		r.database.accountsMU.RLock()
	})
	defer r.database.transactionsMU.RUnlock()
	defer r.database.accountsMU.RUnlock()

	transactions := filter(values(r.database.transactions), func(it beans.Transaction) bool {
		if account, ok := r.database.accounts[it.AccountID]; ok {
			if account.BudgetID == budgetID {
				return true
			}
		}
		return false
	})
	if !from.Empty() {
		transactions = filter(transactions, func(it beans.Transaction) bool { return it.Date.Equal(from.Time) || it.Date.After(from.Time) })
	}
	if !to.Empty() {
		transactions = filter(transactions, func(it beans.Transaction) bool { return it.Date.Equal(to.Time) || it.Date.Before(to.Time) })
	}

	transactions = filter(transactions, func(it beans.Transaction) bool { return !it.CategoryID.Empty() })

	activityByCategory := make(map[beans.ID]beans.Amount)
	for _, t := range transactions {
		if current, ok := activityByCategory[t.CategoryID]; ok {
			sum, err := beans.Arithmetic.Add(current, t.Amount)
			if err != nil {
				panic(err)
			}
			activityByCategory[t.CategoryID] = sum
		} else {
			activityByCategory[t.CategoryID] = t.Amount
		}
	}

	return activityByCategory, nil
}

func (r *transactionRepository) GetIncomeBetween(ctx context.Context, budgetID beans.ID, begin beans.Date, end beans.Date) (beans.Amount, error) {
	r.acquire(func() {
		r.database.transactionsMU.RLock()
		r.database.accountsMU.RLock()
		r.database.categoriesMU.RLock()
	})
	defer r.database.transactionsMU.RUnlock()
	defer r.database.accountsMU.RUnlock()
	defer r.database.categoriesMU.RUnlock()

	// filter to in budget
	transactions := filter(values(r.database.transactions), func(it beans.Transaction) bool {
		if account, ok := r.database.accounts[it.AccountID]; ok {
			if account.BudgetID == budgetID {
				return true
			}
		}
		return false
	})

	// filter by date
	if !begin.Empty() {
		transactions = filter(transactions, func(it beans.Transaction) bool { return it.Date.Equal(begin.Time) || it.Date.After(begin.Time) })
	}
	if !end.Empty() {
		transactions = filter(transactions, func(it beans.Transaction) bool { return it.Date.Equal(end.Time) || it.Date.Before(end.Time) })
	}

	// must have a category
	transactions = filter(transactions, func(it beans.Transaction) bool { return !it.CategoryID.Empty() })

	// category must be income
	transactions = filter(transactions, func(it beans.Transaction) bool {
		if category, ok := r.database.categories[it.CategoryID]; ok {
			if group, ok := r.database.categoryGroups[category.GroupID]; ok {
				return group.IsIncome
			}
		}
		return false
	})

	// sum transactions
	return reduce(transactions, beans.NewAmount(0, 0), func(it beans.Transaction, acc beans.Amount) beans.Amount {
		r, err := beans.Arithmetic.Add(acc, it.Amount)
		if err != nil {
			panic(err)
		}
		return r
	}), nil
}

// helper

func (r *transactionRepository) mapWithRelations(it beans.Transaction) beans.TransactionWithRelations {
	t := beans.TransactionWithRelations{
		Transaction: it,
	}

	account := r.database.accounts[it.AccountID]
	t.Account = beans.RelatedAccount{ID: account.ID, Name: account.Name, OffBudget: account.OffBudget}

	if !it.TransferID.Empty() {
		transfer := r.database.transactions[it.TransferID]
		transferAccount := r.database.accounts[transfer.AccountID]

		t.TransferAccount = beans.OptionalWrap(beans.RelatedAccount{ID: transferAccount.ID, Name: transferAccount.Name, OffBudget: transferAccount.OffBudget})
	}

	t.Variant = beans.GetTransactionVariant(t.Account, t.TransferAccount)

	if !it.CategoryID.Empty() {
		category := r.database.categories[it.CategoryID]
		t.Category = beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name})
	}

	if !it.PayeeID.Empty() {
		payee := r.database.payees[it.PayeeID]
		t.Payee = beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name})
	}
	return t
}
