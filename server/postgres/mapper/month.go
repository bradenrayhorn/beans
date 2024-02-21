package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func Month(d db.Month) (beans.Month, error) {
	budgetID, err := beans.IDFromString(d.BudgetID)
	if err != nil {
		return beans.Month{}, err
	}

	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.Month{}, err
	}

	carryover, err := NumericToAmount(d.Carryover)
	if err != nil {
		return beans.Month{}, err
	}

	return beans.Month{
		ID:        id,
		BudgetID:  budgetID,
		Date:      PgToMonthDate(d.Date),
		Carryover: carryover,
	}, nil
}
