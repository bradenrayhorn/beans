package contract

import "github.com/bradenrayhorn/beans/server/beans"

type contract struct {
	datasource        beans.DataSource
	sessionRepository beans.SessionRepository
}

func (c *contract) ds() beans.DataSource {
	return c.datasource
}

type Contracts struct {
	Account     beans.AccountContract
	Budget      beans.BudgetContract
	Category    beans.CategoryContract
	Month       beans.MonthContract
	Payee       beans.PayeeContract
	Transaction beans.TransactionContract
	User        beans.UserContract
}

func NewContracts(datasource beans.DataSource, sessionRepository beans.SessionRepository) *Contracts {
	contract := contract{datasource, sessionRepository}

	return &Contracts{
		Account:     &accountContract{contract},
		Budget:      &budgetContract{contract},
		Category:    &categoryContract{contract},
		Month:       &monthContract{contract},
		Payee:       &payeeContract{contract},
		Transaction: &transactionContract{contract},
		User:        &userContract{contract},
	}
}
