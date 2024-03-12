package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func Account(d db.Account) (beans.Account, error) {
	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.Account{}, err
	}

	budgetID, err := beans.IDFromString(d.BudgetID)
	if err != nil {
		return beans.Account{}, err
	}

	return beans.Account{
		ID:        id,
		BudgetID:  budgetID,
		Name:      beans.Name(d.Name),
		OffBudget: d.OffBudget,
	}, nil
}

func AccountWithBalance(d db.GetAccountsWithBalanceRow) (beans.AccountWithBalance, error) {
	account, err := Account(d.Account)
	if err != nil {
		return beans.AccountWithBalance{}, err
	}

	balance, err := NumericToAmount(d.Balance)
	if err != nil {
		return beans.AccountWithBalance{}, err
	}
	if balance.Empty() {
		balance = beans.NewAmount(0, 0)
	}

	return beans.AccountWithBalance{
		Account: account,
		Balance: balance,
	}, nil
}
