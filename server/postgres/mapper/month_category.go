package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func MonthCategory(res db.MonthCategory) (beans.MonthCategory, error) {
	id, err := beans.IDFromString(res.ID)
	if err != nil {
		return beans.MonthCategory{}, err
	}
	categoryID, err := beans.IDFromString(res.CategoryID)
	if err != nil {
		return beans.MonthCategory{}, err
	}
	monthID, err := beans.IDFromString(res.MonthID)
	if err != nil {
		return beans.MonthCategory{}, err
	}
	amount, err := NumericToAmount(res.Amount)
	if err != nil {
		return beans.MonthCategory{}, err
	}
	if amount.Empty() {
		amount = beans.NewAmount(0, 0)
	}

	return beans.MonthCategory{
		ID:         id,
		MonthID:    monthID,
		CategoryID: categoryID,
		Amount:     amount,
	}, nil
}

func GetPastMonthCategoriesAssigned(res db.GetPastMonthCategoriesAssignedRow) (beans.ID, beans.Amount, error) {
	id, err := beans.IDFromString(res.CategoryID)
	if err != nil {
		return beans.EmptyID(), beans.NewEmptyAmount(), err
	}
	amount, err := NumericToAmount(res.Assigned)
	if err != nil {
		return beans.EmptyID(), beans.NewEmptyAmount(), err
	}
	if amount.Empty() {
		amount = beans.NewAmount(0, 0)
	}

	return id, amount, err
}
