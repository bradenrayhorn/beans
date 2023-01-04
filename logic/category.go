package logic

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type categoryService struct {
	categoryRepository beans.CategoryRepository
}

func NewCategoryService(categoryRepository beans.CategoryRepository) *categoryService {
	return &categoryService{categoryRepository}
}

func (s categoryService) CreateCategory(ctx context.Context, budgetID beans.ID, groupID beans.ID, name beans.Name) (*beans.Category, error) {
	if err := beans.ValidateFields(
		beans.Field("Budget ID", beans.Required(budgetID)),
		beans.Field("Group ID", beans.Required(groupID)),
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	groupExists, err := s.categoryRepository.GroupExists(ctx, budgetID, groupID)
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

	if err := s.categoryRepository.Create(ctx, nil, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s categoryService) CreateGroup(ctx context.Context, budgetID beans.ID, name beans.Name) (*beans.CategoryGroup, error) {
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

	if err := s.categoryRepository.CreateGroup(ctx, nil, group); err != nil {
		return nil, err
	}

	return group, nil
}
