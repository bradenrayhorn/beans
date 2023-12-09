package beans

import "context"

type Category struct {
	ID       ID
	BudgetID ID
	GroupID  ID
	Name     Name
}

type CategoryGroup struct {
	ID       ID
	BudgetID ID
	Name     Name
	IsIncome bool
}

type CategoryContract interface {
	// Creates a category.
	CreateCategory(ctx context.Context, auth *BudgetAuthContext, groupID ID, name Name) (*Category, error)

	// Creates a category group.
	CreateGroup(ctx context.Context, auth *BudgetAuthContext, name Name) (*CategoryGroup, error)

	// Gets all categories and groups for a budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]*CategoryGroup, []*Category, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, tx Tx, category *Category) error
	GetSingleForBudget(ctx context.Context, id ID, budgetID ID) (*Category, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]*Category, error)
	CreateGroup(ctx context.Context, tx Tx, categoryGroup *CategoryGroup) error
	GetGroupsForBudget(ctx context.Context, budgetID ID) ([]*CategoryGroup, error)
	GroupExists(ctx context.Context, budgetID ID, id ID) (bool, error)
}