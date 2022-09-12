package beans

import "context"

type MonthCategory struct {
	ID         ID
	MonthID    ID
	CategoryID ID
	Amount     Amount
}

type MonthCategoryRepository interface {
	Create(ctx context.Context, monthCategory *MonthCategory) error
	GetForMonth(ctx context.Context, monthID ID) ([]*MonthCategory, error)
}
