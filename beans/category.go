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
}

type CategoryRepository interface {
	Create(context.Context, *Category) error
	GetForBudget(ctx context.Context, budgetID ID) ([]*Category, error)
	CreateGroup(context.Context, *CategoryGroup) error
	GetGroupsForBudget(ctx context.Context, budgetID ID) ([]*CategoryGroup, error)
	GroupExists(ctx context.Context, budgetID ID, id ID) (bool, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, budgetID ID, groupID ID, name Name) (*Category, error)
	CreateGroup(ctx context.Context, budgetID ID, name Name) (*CategoryGroup, error)
}
