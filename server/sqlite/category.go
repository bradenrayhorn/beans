package sqlite

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"golang.org/x/net/context"
	"zombiezen.com/go/sqlite"
)

type categoryRepository struct{ repository }

var _ beans.CategoryRepository = (*categoryRepository)(nil)

const categoryCreateSQL = `
INSERT INTO categories (id, budget_id, group_id, name)
	VALUES (:id, :budgetID, :groupID, :name)
`

func (r *categoryRepository) Create(ctx context.Context, tx beans.Tx, category beans.Category) error {
	return db[any](r.pool).
		inTx(tx).
		execute(ctx, categoryCreateSQL, map[string]any{
			":id":       category.ID.String(),
			":budgetID": category.BudgetID.String(),
			":groupID":  category.GroupID.String(),
			":name":     string(category.Name),
		})
}

const getCategorySQL = `
SELECT * FROM categories WHERE id = :id AND budget_id = :budgetID
`

func (r *categoryRepository) GetSingleForBudget(ctx context.Context, id beans.ID, budgetID beans.ID) (beans.Category, error) {
	return db[beans.Category](r.pool).
		mapWith(mapCategory).
		one(ctx, getCategorySQL, map[string]any{
			":id":       id.String(),
			":budgetID": budgetID.String(),
		})
}

const getCategoryGroupSQL = `
SELECT * FROM category_groups WHERE id = :id AND budget_id = :budgetID
`

func (r *categoryRepository) GetCategoryGroup(ctx context.Context, id beans.ID, budgetID beans.ID) (beans.CategoryGroup, error) {
	return db[beans.CategoryGroup](r.pool).
		mapWith(mapCategoryGroup).
		one(ctx, getCategoryGroupSQL, map[string]any{
			":id":       id.String(),
			":budgetID": budgetID.String(),
		})
}

const getCategoriesForGroupSQL = `
SELECT categories.* FROM categories
JOIN category_groups ON category_groups.id = categories.group_id
	AND category_groups.id = :groupID
	AND category_groups.budget_id = :budgetID
`

func (r *categoryRepository) GetCategoriesForGroup(ctx context.Context, id beans.ID, budgetID beans.ID) ([]beans.Category, error) {
	return db[beans.Category](r.pool).
		mapWith(mapCategory).
		many(ctx, getCategoriesForGroupSQL, map[string]any{
			":groupID":  id.String(),
			":budgetID": budgetID.String(),
		})
}

const getCategoriesForBudgetSQL = `
SELECT * FROM categories WHERE budget_id = :budgetID
`

func (r *categoryRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Category, error) {
	return db[beans.Category](r.pool).
		mapWith(mapCategory).
		many(ctx, getCategoriesForBudgetSQL, map[string]any{
			":budgetID": budgetID.String(),
		})
}

const categoryGroupCreateSQL = `
INSERT INTO category_groups (id, budget_id, name, is_income)
	VALUES (:id, :budgetID, :name, :isIncome)
`

func (r *categoryRepository) CreateGroup(ctx context.Context, tx beans.Tx, category beans.CategoryGroup) error {
	return db[any](r.pool).
		inTx(tx).
		execute(ctx, categoryGroupCreateSQL, map[string]any{
			":id":       category.ID.String(),
			":budgetID": category.BudgetID.String(),
			":name":     string(category.Name),
			":isIncome": category.IsIncome,
		})
}

const getCategoryGroupsForBudgetSQL = `
SELECT * FROM category_groups WHERE budget_id = :budgetID
`

func (r *categoryRepository) GetGroupsForBudget(ctx context.Context, budgetID beans.ID) ([]beans.CategoryGroup, error) {
	return db[beans.CategoryGroup](r.pool).
		mapWith(mapCategoryGroup).
		many(ctx, getCategoryGroupsForBudgetSQL, map[string]any{
			":budgetID": budgetID.String(),
		})
}

// mappers

func mapCategory(stmt *sqlite.Stmt) (beans.Category, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Category{}, err
	}
	budgetID, err := mapID(stmt, "budget_id")
	if err != nil {
		return beans.Category{}, err
	}
	groupID, err := mapID(stmt, "group_id")
	if err != nil {
		return beans.Category{}, err
	}

	return beans.Category{
		ID:       id,
		BudgetID: budgetID,
		GroupID:  groupID,
		Name:     beans.Name(stmt.GetText("name")),
	}, nil
}

func mapCategoryGroup(stmt *sqlite.Stmt) (beans.CategoryGroup, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.CategoryGroup{}, err
	}
	budgetID, err := mapID(stmt, "budget_id")
	if err != nil {
		return beans.CategoryGroup{}, err
	}

	return beans.CategoryGroup{
		ID:       id,
		BudgetID: budgetID,
		IsIncome: stmt.GetBool("is_income"),
		Name:     beans.Name(stmt.GetText("name")),
	}, nil
}
