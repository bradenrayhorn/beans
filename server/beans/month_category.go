package beans

import "context"

type MonthCategory struct {
	ID         ID
	MonthID    ID
	CategoryID ID
	Amount     Amount
}

type MonthCategoryWithDetails struct {
	ID         ID
	CategoryID ID
	Amount     Amount
	Activity   Amount
	Available  Amount
}

type MonthCategoryRepository interface {
	Create(ctx context.Context, tx Tx, monthCategory MonthCategory) error
	UpdateAmount(ctx context.Context, monthCategory MonthCategory) error
	GetForMonth(ctx context.Context, month Month) ([]MonthCategoryWithDetails, error)

	// Gets the month category, or creates it if it does not exist.
	GetOrCreate(ctx context.Context, tx Tx, month Month, categoryID ID) (MonthCategory, error)

	// Gets the amount assigned in a month.
	GetAssignedInMonth(ctx context.Context, month Month) (Amount, error)
}
