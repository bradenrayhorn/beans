package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type monthContract struct{ contract }

var _ beans.MonthContract = (*monthContract)(nil)

func (c *monthContract) GetOrCreate(ctx context.Context, auth *beans.BudgetAuthContext, date beans.MonthDate) (beans.MonthWithDetails, error) {
	month, err := c.createMonth(ctx, auth, date)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	pastMonth, err := c.ds().MonthRepository().GetOrCreate(
		ctx,
		nil,
		auth.BudgetID(),
		month.Date.Previous(),
	)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	categories, err := c.services.MonthCategory.GetForMonth(ctx, month)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	income, err := c.ds().TransactionRepository().GetIncomeBetween(ctx, auth.BudgetID(), month.Date.FirstDay(), month.Date.LastDay())
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	assignedInMonth, err := c.ds().MonthCategoryRepository().GetAssignedInMonth(ctx, month)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	available, err := beans.Arithmetic.Add(
		income,
		pastMonth.Carryover,
		beans.Arithmetic.Negate(month.Carryover),
		beans.Arithmetic.Negate(assignedInMonth),
	)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	return beans.MonthWithDetails{
		Month: month,

		CarriedOver: pastMonth.Carryover,
		Income:      income,
		Assigned:    assignedInMonth,
		Budgetable:  available,
		Categories:  categories,
	}, nil
}

func (c *monthContract) createMonth(ctx context.Context, auth *beans.BudgetAuthContext, date beans.MonthDate) (beans.Month, error) {
	return beans.ExecTx(ctx, c.ds().TxManager(), func(tx beans.Tx) (beans.Month, error) {
		// create month
		month, err := c.ds().MonthRepository().GetOrCreate(ctx, tx, auth.BudgetID(), date)
		if err != nil {
			return beans.Month{}, err
		}

		// create month categories for every category
		categories, err := c.ds().CategoryRepository().GetForBudget(ctx, auth.BudgetID())
		if err != nil {
			return beans.Month{}, err
		}
		for _, category := range categories {
			if _, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, tx, month, category.ID); err != nil {
				return beans.Month{}, err
			}
		}

		return month, nil
	})
}
func (c *monthContract) Update(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID, carryover beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Carryover", beans.Required(&carryover), beans.Positive(carryover)),
	); err != nil {
		return err
	}

	month, err := c.ds().MonthRepository().Get(ctx, auth.BudgetID(), monthID)
	if err != nil {
		return err
	}

	month.Carryover = carryover

	return c.ds().MonthRepository().Update(ctx, month)
}

func (c *monthContract) SetCategoryAmount(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID, categoryID beans.ID, amount beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Amount", beans.Required(&amount), beans.NonZero(amount), beans.Positive(amount)),
	); err != nil {
		return err
	}

	month, err := c.ds().MonthRepository().Get(ctx, auth.BudgetID(), monthID)
	if err != nil {
		return err
	}

	if _, err := c.ds().CategoryRepository().GetSingleForBudget(ctx, categoryID, auth.BudgetID()); err != nil {
		return err
	}

	monthCategory, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, nil, month, categoryID)
	if err != nil {
		return err
	}

	monthCategory.Amount = amount

	return c.ds().MonthCategoryRepository().UpdateAmount(ctx, monthCategory)
}
