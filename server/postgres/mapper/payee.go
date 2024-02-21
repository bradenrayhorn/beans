package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func Payee(d db.Payee) (beans.Payee, error) {
	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.Payee{}, err
	}

	budgetID, err := beans.IDFromString(d.BudgetID)
	if err != nil {
		return beans.Payee{}, err
	}

	return beans.Payee{
		ID:       id,
		BudgetID: budgetID,
		Name:     beans.Name(d.Name),
	}, nil
}
