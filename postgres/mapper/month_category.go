package mapper

import (
	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
)

func GetMonthCategoriesForMonthRow(res db.GetMonthCategoriesForMonthRow) (*beans.MonthCategory, error) {
	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return nil, err
	}
	monthID, err := beans.BeansIDFromString(res.MonthID)
	if err != nil {
		return nil, err
	}
	categoryID, err := beans.BeansIDFromString(res.CategoryID)
	if err != nil {
		return nil, err
	}
	amount, err := NumericToAmount(res.Amount)
	if err != nil {
		return nil, err
	}
	if amount.Empty() {
		amount = beans.NewAmount(0, 0)
	}

	activity, err := NumericToAmount(res.Activity)
	if err != nil {
		return nil, err
	}
	if activity.Empty() {
		activity = beans.NewAmount(0, 0)
	}

	return &beans.MonthCategory{
		ID:         id,
		MonthID:    monthID,
		CategoryID: categoryID,
		Amount:     amount,

		Activity: activity,
	}, nil
}
