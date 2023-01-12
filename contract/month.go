package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type monthContract struct {
	monthRepository         beans.MonthRepository
	monthCategoryRepository beans.MonthCategoryRepository
}

func NewMonthContract(
	monthRepository beans.MonthRepository,
	monthCategoryRepository beans.MonthCategoryRepository,
) *monthContract {
	return &monthContract{monthRepository, monthCategoryRepository}
}

func (c *monthContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID) (*beans.Month, []*beans.MonthCategory, error) {
	month, err := c.getAndVerifyMonth(ctx, auth, monthID)
	if err != nil {
		return nil, nil, err
	}

	categories, err := c.monthCategoryRepository.GetForMonth(ctx, month)
	if err != nil {
		return nil, nil, err
	}

	return month, categories, nil
}

func (c *monthContract) CreateMonth(ctx context.Context, auth *beans.BudgetAuthContext, date beans.MonthDate) (*beans.Month, error) {
	return c.monthRepository.GetOrCreate(ctx, auth.BudgetID(), date)
}

func (c *monthContract) SetCategoryAmount(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID, categoryID beans.ID, amount beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Amount", beans.NonZero(amount), beans.Positive(amount)),
	); err != nil {
		return err
	}

	_, err := c.getAndVerifyMonth(ctx, auth, monthID)
	if err != nil {
		return err
	}

	monthCategory, err := c.monthCategoryRepository.GetOrCreate(ctx, monthID, categoryID)
	if err != nil {
		return err
	}

	return c.monthCategoryRepository.UpdateAmount(ctx, monthCategory.ID, amount)
}

func (c *monthContract) getAndVerifyMonth(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID) (*beans.Month, error) {
	month, err := c.monthRepository.Get(ctx, monthID)
	if err != nil {
		return nil, err
	}

	if month.BudgetID != auth.BudgetID() {
		return nil, beans.NewError(beans.EFORBIDDEN, "No access to month")
	}

	return month, nil
}
