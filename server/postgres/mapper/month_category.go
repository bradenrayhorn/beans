package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func GetMonthCategoriesForMonthRow(res db.GetMonthCategoriesForMonthRow) (beans.MonthCategoryWithDetails, error) {
	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return beans.MonthCategoryWithDetails{}, err
	}
	categoryID, err := beans.BeansIDFromString(res.CategoryID)
	if err != nil {
		return beans.MonthCategoryWithDetails{}, err
	}
	amount, err := NumericToAmount(res.Amount)
	if err != nil {
		return beans.MonthCategoryWithDetails{}, err
	}
	if amount.Empty() {
		amount = beans.NewAmount(0, 0)
	}

	activity, err := NumericToAmount(res.Activity)
	if err != nil {
		return beans.MonthCategoryWithDetails{}, err
	}
	if activity.Empty() {
		activity = beans.NewAmount(0, 0)
	}

	return beans.MonthCategoryWithDetails{
		ID:         id,
		CategoryID: categoryID,
		Amount:     amount,

		Activity: activity,
	}, nil
}
