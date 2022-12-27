package beans

import "context"

type MonthCategory struct {
	ID         ID
	MonthID    ID
	CategoryID ID
	Amount     Amount

	// Must be explicitly loaded.
	Spent Amount
}

type MonthCategoryService interface {
	CreateOrUpdate(ctx context.Context, monthID ID, categoryID ID, amount Amount) error
	CreateIfNotExists(ctx context.Context, monthID ID, categoryID ID) error
}

type MonthCategoryRepository interface {
	Create(ctx context.Context, monthCategory *MonthCategory) error
	UpdateAmount(ctx context.Context, monthCategoryID ID, amount Amount) error
	// Gets categories by month. Attached Spent field.
	GetForMonth(ctx context.Context, month Month) ([]*MonthCategory, error)
	GetByMonthAndCategory(ctx context.Context, monthID ID, categoryID ID) (*MonthCategory, error)
}
