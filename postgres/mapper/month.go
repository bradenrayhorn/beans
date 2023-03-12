package mapper

import (
	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
)

func Month(d db.Month) (*beans.Month, error) {
	budgetID, err := beans.BeansIDFromString(d.BudgetID)
	if err != nil {
		return nil, err
	}

	id, err := beans.BeansIDFromString(d.ID)
	if err != nil {
		return nil, err
	}

	carryover, err := NumericToAmount(d.Carryover)
	if err != nil {
		return nil, err
	}

	return &beans.Month{
		ID:        id,
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(beans.NewDate(d.Date)),
		Carryover: carryover,
	}, nil
}
