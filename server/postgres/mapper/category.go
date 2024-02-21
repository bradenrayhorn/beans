package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

func CategoryGroup(d db.CategoryGroup) (beans.CategoryGroup, error) {
	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.CategoryGroup{}, err
	}

	budgetID, err := beans.IDFromString(d.BudgetID)
	if err != nil {
		return beans.CategoryGroup{}, err
	}

	return beans.CategoryGroup{
		ID:       id,
		BudgetID: budgetID,
		Name:     beans.Name(d.Name),
		IsIncome: d.IsIncome,
	}, nil
}

func Category(d db.Category) (beans.Category, error) {
	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.Category{}, err
	}

	budgetID, err := beans.IDFromString(d.BudgetID)
	if err != nil {
		return beans.Category{}, err
	}

	groupID, err := beans.IDFromString(d.GroupID)
	if err != nil {
		return beans.Category{}, err
	}

	return beans.Category{
		ID:       id,
		BudgetID: budgetID,
		GroupID:  groupID,
		Name:     beans.Name(d.Name),
	}, nil
}
