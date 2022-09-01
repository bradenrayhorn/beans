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
}
