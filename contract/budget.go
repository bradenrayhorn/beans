package contract

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type budgetContract struct {
	budgetRepository   beans.BudgetRepository
	categoryRepository beans.CategoryRepository
	monthRepository    beans.MonthRepository
	txManager          beans.TxManager
}

func NewBudgetContract(
	budgetRepository beans.BudgetRepository,
	categoryRepository beans.CategoryRepository,
	monthRepository beans.MonthRepository,
	txManager beans.TxManager,
) *budgetContract {
	return &budgetContract{budgetRepository, categoryRepository, monthRepository, txManager}
}

func (c *budgetContract) Create(ctx context.Context, name beans.Name, userID beans.UserID) (*beans.Budget, error) {
	if err := beans.ValidateFields(beans.Field("Budget name", name)); err != nil {
		return nil, err
	}

	budgetID := beans.NewBeansID()

	tx, err := c.txManager.Create(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// create budget
	if err := c.budgetRepository.Create(ctx, tx, budgetID, name, userID); err != nil {
		return nil, err
	}

	// create month
	if err := c.monthRepository.Create(ctx, tx, &beans.Month{
		ID:       beans.NewBeansID(),
		BudgetID: budgetID,
		Date:     beans.NewDate(beans.NormalizeMonth(time.Now())),
	}); err != nil {
		return nil, err
	}

	// create income category
	categoryGroup := &beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		BudgetID: budgetID,
		Name:     "Income",
	}
	if err := c.categoryRepository.CreateGroup(ctx, tx, categoryGroup); err != nil {
		return nil, err
	}

	category := &beans.Category{
		ID:       beans.NewBeansID(),
		GroupID:  categoryGroup.ID,
		BudgetID: budgetID,
		Name:     "Income",
		IsIncome: true,
	}
	if err := c.categoryRepository.Create(ctx, tx, category); err != nil {
		return nil, err
	}

	// commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &beans.Budget{
		ID:   budgetID,
		Name: name,
	}, nil
}

func (c *budgetContract) Get(ctx context.Context, id beans.ID, userID beans.UserID) (*beans.Budget, *beans.Month, error) {
	budget, err := c.budgetRepository.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if !budget.UserHasAccess(userID) {
		return nil, nil, beans.ErrorNotFound
	}

	month, err := c.monthRepository.GetLatest(ctx, budget.ID)
	if err != nil {
		return nil, nil, beans.WrapError(err, beans.ErrorInternal)
	}

	return budget, month, nil
}

func (c *budgetContract) GetAll(ctx context.Context, userID beans.UserID) ([]*beans.Budget, error) {
	return c.budgetRepository.GetBudgetsForUser(ctx, userID)
}
