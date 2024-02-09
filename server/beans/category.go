package beans

import "context"

type Category struct {
	ID       ID
	BudgetID ID
	GroupID  ID
	Name     Name
}

type RelatedCategory struct {
	ID   ID
	Name Name
}

type CategoryGroup struct {
	ID       ID
	BudgetID ID
	Name     Name
	IsIncome bool
}

type CategoryGroupWithCategories struct {
	CategoryGroup
	Categories []Category
}

type CategoryContract interface {
	// Creates a category.
	CreateCategory(ctx context.Context, auth *BudgetAuthContext, groupID ID, name Name) (Category, error)

	// Creates a category group.
	CreateGroup(ctx context.Context, auth *BudgetAuthContext, name Name) (CategoryGroup, error)

	// Gets all categories and groups for a budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]CategoryGroupWithCategories, error)

	// Gets group for a budget.
	GetGroup(ctx context.Context, auth *BudgetAuthContext, id ID) (CategoryGroupWithCategories, error)

	// Gets category for a budget.
	GetCategory(ctx context.Context, auth *BudgetAuthContext, id ID) (Category, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, tx Tx, category Category) error
	GetSingleForBudget(ctx context.Context, id ID, budgetID ID) (Category, error)
	GetCategoryGroup(ctx context.Context, id ID, budgetID ID) (CategoryGroup, error)
	GetCategoriesForGroup(ctx context.Context, id ID, budgetID ID) ([]Category, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]Category, error)
	CreateGroup(ctx context.Context, tx Tx, categoryGroup CategoryGroup) error
	GetGroupsForBudget(ctx context.Context, budgetID ID) ([]CategoryGroup, error)
	GroupExists(ctx context.Context, budgetID ID, id ID) (bool, error)
}
