package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type categoryContract struct{ contract }

var _ beans.CategoryContract = (*categoryContract)(nil)

func (c *categoryContract) CreateCategory(ctx context.Context, auth *beans.BudgetAuthContext, groupID beans.ID, name beans.Name) (beans.Category, error) {
	if err := beans.ValidateFields(
		beans.Field("Group ID", beans.Required(groupID)),
		beans.Field("Name", name),
	); err != nil {
		return beans.Category{}, err
	}

	category := beans.Category{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		GroupID:  groupID,
		Name:     name,
	}

	err := beans.ExecTxNil(ctx, c.ds().TxManager(), func(tx beans.Tx) error {
		_, err := c.ds().CategoryRepository().GetCategoryGroup(ctx, groupID, auth.BudgetID())
		if err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.NewError(beans.EINVALID, "Invalid Group ID.")
			}
			return err
		}

		if err := c.ds().CategoryRepository().Create(ctx, nil, category); err != nil {
			return err
		}

		// create month categories for existing months
		months, err := c.ds().MonthRepository().GetForBudget(ctx, auth.BudgetID())
		if err != nil {
			return err
		}

		for _, month := range months {
			err = c.ds().MonthCategoryRepository().Create(ctx, tx, beans.MonthCategory{
				ID:         beans.NewBeansID(),
				MonthID:    month.ID,
				CategoryID: category.ID,
				Amount:     beans.NewAmount(0, 0),
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return beans.Category{}, err
	}

	return category, nil
}

func (c *categoryContract) CreateGroup(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (beans.CategoryGroup, error) {
	if err := beans.ValidateFields(
		beans.Field("Name", name),
	); err != nil {
		return beans.CategoryGroup{}, err
	}

	group := beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		Name:     name,
	}

	if err := c.ds().CategoryRepository().CreateGroup(ctx, nil, group); err != nil {
		return beans.CategoryGroup{}, err
	}

	return group, nil
}

func (c *categoryContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.CategoryGroupWithCategories, error) {
	groups, err := c.ds().CategoryRepository().GetGroupsForBudget(ctx, auth.BudgetID())
	if err != nil {
		return nil, err
	}

	categories, err := c.ds().CategoryRepository().GetForBudget(ctx, auth.BudgetID())
	if err != nil {
		return nil, err
	}

	// group categories by group
	categoriesByGroup := make(map[string][]beans.Category)
	for _, group := range groups {
		categoriesByGroup[group.ID.String()] = make([]beans.Category, 0)
	}
	for _, category := range categories {
		groupID := category.GroupID.String()
		categoriesByGroup[groupID] = append(categoriesByGroup[groupID], category)
	}

	// associate categories with their groups
	groupsWithCategories := make([]beans.CategoryGroupWithCategories, len(groups))
	for i, group := range groups {
		groupsWithCategories[i] = beans.CategoryGroupWithCategories{
			CategoryGroup: group,
			Categories:    categoriesByGroup[group.ID.String()],
		}
	}

	return groupsWithCategories, nil
}

func (c *categoryContract) GetGroup(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.CategoryGroupWithCategories, error) {
	group, err := c.ds().CategoryRepository().GetCategoryGroup(ctx, id, auth.BudgetID())
	if err != nil {
		return beans.CategoryGroupWithCategories{}, err
	}

	categories, err := c.ds().CategoryRepository().GetCategoriesForGroup(ctx, id, auth.BudgetID())
	if err != nil {
		return beans.CategoryGroupWithCategories{}, err
	}

	return beans.CategoryGroupWithCategories{
		CategoryGroup: group,
		Categories:    categories,
	}, nil
}

func (c *categoryContract) GetCategory(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.Category, error) {
	category, err := c.ds().CategoryRepository().GetSingleForBudget(ctx, id, auth.BudgetID())
	if err != nil {
		return beans.Category{}, err
	}

	return category, nil
}
