package beans

import "context"

type MonthCategory struct {
	ID         ID
	MonthID    ID
	CategoryID ID
	Amount     Amount
}

type MonthCategoryService interface {
	CreateOrUpdate(ctx context.Context, monthID ID, categoryID ID, amount Amount) error
}

type MonthCategoryRepository interface {
	Create(ctx context.Context, monthCategory *MonthCategory) error
	UpdateAmount(ctx context.Context, monthCategoryID ID, amount Amount) error
	GetForMonth(ctx context.Context, monthID ID) ([]*MonthCategory, error)
	GetByMonthAndCategory(ctx context.Context, monthID ID, categoryID ID) (*MonthCategory, error)
}
