package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type categoryContract struct {
	categoryRepository      beans.CategoryRepository
	monthCategoryRepository beans.MonthCategoryRepository
	monthRepository         beans.MonthRepository
	txManager               beans.TxManager
}

func NewCategoryContract(
	categoryRepository beans.CategoryRepository,
	monthCategoryRepository beans.MonthCategoryRepository,
	monthRepository beans.MonthRepository,
	txManager beans.TxManager,
) *categoryContract {
	return &categoryContract{
		categoryRepository,
		monthCategoryRepository,
		monthRepository,
		txManager,
	}
}

func (c *categoryContract) CreateCategory(ctx context.Context, auth *beans.BudgetAuthContext, groupID beans.ID, name beans.Name) (*beans.Category, error) {
	if err := beans.ValidateFields(
		beans.Field("Group ID", beans.Required(groupID)),
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	tx, err := c.txManager.Create(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	groupExists, err := c.categoryRepository.GroupExists(ctx, auth.BudgetID(), groupID)
	if err != nil {
		return nil, err
	}
	if !groupExists {
		return nil, beans.NewError(beans.EINVALID, "Invalid Group ID.")
	}

	category := &beans.Category{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		GroupID:  groupID,
		Name:     name,
	}

	if err := c.categoryRepository.Create(ctx, nil, category); err != nil {
		return nil, err
	}

	// create month categories for existing months
	months, err := c.monthRepository.GetForBudget(ctx, auth.BudgetID())
	if err != nil {
		return nil, err
	}

	for _, month := range months {
		err = c.monthCategoryRepository.Create(ctx, tx, &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    month.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(0, 0),
		})

		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryContract) CreateGroup(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (*beans.CategoryGroup, error) {
	if err := beans.ValidateFields(
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	group := &beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		Name:     name,
	}

	if err := c.categoryRepository.CreateGroup(ctx, nil, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (c *categoryContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]*beans.CategoryGroup, []*beans.Category, error) {
	groups, err := c.categoryRepository.GetGroupsForBudget(ctx, auth.BudgetID())
	if err != nil {
		return nil, nil, err
	}

	categories, err := c.categoryRepository.GetForBudget(ctx, auth.BudgetID())
	if err != nil {
		return nil, nil, err
	}

	return groups, categories, nil
}
