package mapper

import (
	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
)

func Payee(d db.Payee) (*beans.Payee, error) {
	id, err := beans.BeansIDFromString(d.ID)
	if err != nil {
		return nil, err
	}

	budgetID, err := beans.BeansIDFromString(d.BudgetID)
	if err != nil {
		return nil, err
	}

	return &beans.Payee{
		ID:       id,
		BudgetID: budgetID,
		Name:     beans.Name(d.Name),
	}, nil
}
