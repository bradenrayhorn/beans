package postgres

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
	"golang.org/x/net/context"
)

type categoryRepository struct {
	repository
}

func NewCategoryRepository(pool *DbPool) *categoryRepository {
	return &categoryRepository{repository{pool}}
}

func (r *categoryRepository) Create(ctx context.Context, tx beans.Tx, category *beans.Category) error {
	return r.DB(tx).CreateCategory(ctx, db.CreateCategoryParams{
		ID:       category.ID.String(),
		Name:     string(category.Name),
		BudgetID: category.BudgetID.String(),
		GroupID:  category.GroupID.String(),
	})
}

func (r *categoryRepository) GetSingleForBudget(ctx context.Context, id beans.ID, budgetID beans.ID) (*beans.Category, error) {
	res, err := r.DB(nil).GetCategoryForBudget(ctx, db.GetCategoryForBudgetParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.Category(res)
}

func (r *categoryRepository) GetCategoryGroup(ctx context.Context, id beans.ID, budgetID beans.ID) (*beans.CategoryGroup, error) {
	res, err := r.DB(nil).GetCategoryGroup(ctx, db.GetCategoryGroupParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.CategoryGroup(res)
}

func (r *categoryRepository) GetCategoriesForGroup(ctx context.Context, id beans.ID, budgetID beans.ID) ([]*beans.Category, error) {
	res, err := r.DB(nil).GetCategoriesForGroup(ctx, db.GetCategoriesForGroupParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.Category)
}

func (r *categoryRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Category, error) {
	res, err := r.DB(nil).GetCategoriesForBudget(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.Category)
}

func (r *categoryRepository) CreateGroup(ctx context.Context, tx beans.Tx, category *beans.CategoryGroup) error {
	return r.DB(tx).CreateCategoryGroup(ctx, db.CreateCategoryGroupParams{
		ID:       category.ID.String(),
		BudgetID: category.BudgetID.String(),
		Name:     string(category.Name),
		IsIncome: category.IsIncome,
	})
}

func (r *categoryRepository) GetGroupsForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.CategoryGroup, error) {
	res, err := r.DB(nil).GetCategoryGroupsForBudget(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.CategoryGroup)
}

func (r *categoryRepository) GroupExists(ctx context.Context, budgetID beans.ID, id beans.ID) (bool, error) {
	return r.DB(nil).CategoryGroupExists(ctx, db.CategoryGroupExistsParams{BudgetID: budgetID.String(), ID: id.String()})
}
