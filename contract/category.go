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

func (c *categoryContract) CreateCategory(ctx context.Context, auth *beans.BudgetAuthContext, groupID beans.ID, name beans.Name) (*beans.Category, error) {
	if err := beans.ValidateFields(
		beans.Field("Budget ID", beans.Required(auth.BudgetID())),
		beans.Field("Group ID", beans.Required(groupID)),
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

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

	return category, nil
}

func (c *categoryContract) CreateGroup(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (*beans.CategoryGroup, error) {
	if err := beans.ValidateFields(
		beans.Field("Budget ID", beans.Required(auth.BudgetID())),
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
