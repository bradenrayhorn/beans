package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type categoryContract struct {
	categoryRepository beans.CategoryRepository
}

func NewCategoryContract(
	categoryRepository beans.CategoryRepository,
) *categoryContract {
	return &categoryContract{categoryRepository}
}

func (c *categoryContract) CreateCategory(ctx context.Context, budgetID beans.ID, groupID beans.ID, name beans.Name) (*beans.Category, error) {
	if err := beans.ValidateFields(
		beans.Field("Budget ID", beans.Required(budgetID)),
		beans.Field("Group ID", beans.Required(groupID)),
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	groupExists, err := c.categoryRepository.GroupExists(ctx, budgetID, groupID)
	if err != nil {
		return nil, err
	}
	if !groupExists {
		return nil, beans.NewError(beans.EINVALID, "Invalid Group ID.")
	}

	category := &beans.Category{
		ID:       beans.NewBeansID(),
		BudgetID: budgetID,
		GroupID:  groupID,
		Name:     name,
	}

	if err := c.categoryRepository.Create(ctx, nil, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryContract) CreateGroup(ctx context.Context, budgetID beans.ID, name beans.Name) (*beans.CategoryGroup, error) {
	if err := beans.ValidateFields(
		beans.Field("Budget ID", beans.Required(budgetID)),
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	group := &beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		BudgetID: budgetID,
		Name:     name,
	}

	if err := c.categoryRepository.CreateGroup(ctx, nil, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (c *categoryContract) GetAll(ctx context.Context, budgetID beans.ID) ([]*beans.CategoryGroup, []*beans.Category, error) {
	groups, err := c.categoryRepository.GetGroupsForBudget(ctx, budgetID)
	if err != nil {
		return nil, nil, err
	}

	categories, err := c.categoryRepository.GetForBudget(ctx, budgetID)
	if err != nil {
		return nil, nil, err
	}

	return groups, categories, nil
}
