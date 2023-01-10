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

type MonthCategoryRepository interface {
	Create(ctx context.Context, monthCategory *MonthCategory) error
	UpdateAmount(ctx context.Context, monthCategoryID ID, amount Amount) error
	// Gets categories by month. Attaches Spent field.
	GetForMonth(ctx context.Context, month *Month) ([]*MonthCategory, error)
	// Gets the month category, or creates it if it does not exist.
	GetOrCreate(ctx context.Context, monthID ID, categoryID ID) (*MonthCategory, error)
}
