package postgres

import (
	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
)

type categoryRepository struct {
	db *db.Queries
}

func NewCategoryRepository(pool *pgxpool.Pool) *categoryRepository {
	return &categoryRepository{db: db.New(pool)}
}

func (r *categoryRepository) Create(ctx context.Context, category *beans.Category) error {
	return r.db.CreateCategory(ctx, db.CreateCategoryParams{
		ID:       category.ID.String(),
		BudgetID: category.BudgetID.String(),
		Name:     string(category.Name),
		GroupID:  category.GroupID.String(),
	})
}

func (r *categoryRepository) GetSingleForBudget(ctx context.Context, id beans.ID, budgetID beans.ID) (*beans.Category, error) {
	dbCategory, err := r.db.GetCategoryForBudget(ctx, db.GetCategoryForBudgetParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	groupID, err := beans.BeansIDFromString(dbCategory.GroupID)
	if err != nil {
		return nil, err
	}
	return &beans.Category{
		ID:       id,
		BudgetID: budgetID,
		GroupID:  groupID,
		Name:     beans.Name(dbCategory.Name),
	}, nil
}

func (r *categoryRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Category, error) {
	var categories []*beans.Category
	dbCategories, err := r.db.GetCategoriesForBudget(ctx, budgetID.String())
	if err != nil {
		return categories, mapPostgresError(err)
	}
	for _, c := range dbCategories {
		id, err := beans.BeansIDFromString(c.ID)
		if err != nil {
			return categories, err
		}
		groupID, err := beans.BeansIDFromString(c.GroupID)
		if err != nil {
			return categories, err
		}
		categories = append(categories, &beans.Category{ID: id, BudgetID: budgetID, Name: beans.Name(c.Name), GroupID: groupID})
	}

	return categories, nil
}

func (r *categoryRepository) CreateGroup(ctx context.Context, category *beans.CategoryGroup) error {
	return r.db.CreateCategoryGroup(ctx, db.CreateCategoryGroupParams{
		ID:       category.ID.String(),
		BudgetID: category.BudgetID.String(),
		Name:     string(category.Name),
	})
}

func (r *categoryRepository) GetGroupsForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.CategoryGroup, error) {
	var categoryGroups []*beans.CategoryGroup
	dbCategoryGroups, err := r.db.GetCategoryGroupsForBudget(ctx, budgetID.String())
	if err != nil {
		return categoryGroups, nil
	}
	for _, c := range dbCategoryGroups {
		id, err := beans.BeansIDFromString(c.ID)
		if err != nil {
			return categoryGroups, nil
		}
		categoryGroups = append(categoryGroups, &beans.CategoryGroup{ID: id, BudgetID: budgetID, Name: beans.Name(c.Name)})
	}

	return categoryGroups, nil
}

func (r *categoryRepository) GroupExists(ctx context.Context, budgetID beans.ID, id beans.ID) (bool, error) {
	return r.db.CategoryGroupExists(ctx, db.CategoryGroupExistsParams{BudgetID: budgetID.String(), ID: id.String()})
}
