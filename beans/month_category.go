package beans

import "context"

type MonthCategory struct {
	ID         ID
	MonthID    ID
	CategoryID ID
	Amount     Amount

	// Must be explicitly loaded.
	Activity Amount
	// Must be explicitly loaded.
	Available Amount
}

type MonthCategoryRepository interface {
	Create(ctx context.Context, tx Tx, monthCategory *MonthCategory) error
	UpdateAmount(ctx context.Context, monthCategoryID ID, amount Amount) error
	// Gets categories by month. Attaches Activity, Available fields.
	GetForMonth(ctx context.Context, month *Month) ([]*MonthCategory, error)
	// Gets the month category, or creates it if it does not exist.
	GetOrCreate(ctx context.Context, tx Tx, monthID ID, categoryID ID) (*MonthCategory, error)

	// Gets the amount assigned in a month.
	GetAssignedInMonth(ctx context.Context, monthID ID) (Amount, error)
}
