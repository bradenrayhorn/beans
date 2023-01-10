package beans

import (
	"time"

	"golang.org/x/net/context"
)

type Month struct {
	ID       ID
	BudgetID ID
	Date     MonthDate
}

type MonthDate struct {
	date Date
}

func (d MonthDate) String() string {
	return d.date.String()
}

func (d MonthDate) Time() time.Time {
	return d.date.Time
}

func NewMonthDate(date Date) MonthDate {
	return MonthDate{date: NewDate(normalizeMonth(date.Time))}
}

func (m Month) String() string {
	return m.Date.date.Format("2006.01")
}

func normalizeMonth(date time.Time) time.Time {
	date = date.AddDate(0, 0, -date.Day()+1)

	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

type MonthContract interface {
	// Gets a month and its categories.
	Get(ctx context.Context, monthID ID) (*Month, []*MonthCategory, error)

	// Creates a month in the budget with the provided date. If the month already exists, it is returned instead with no error.
	CreateMonth(ctx context.Context, budgetID ID, date MonthDate) (*Month, error)

	// Sets the assigned amount on a category for a month. The user must have access to the month.
	SetCategoryAmount(ctx context.Context, monthID ID, categoryID ID, amount Amount) error
}

type MonthRepository interface {
	Create(ctx context.Context, tx Tx, month *Month) error
	Get(ctx context.Context, id ID) (*Month, error)
	GetOrCreate(ctx context.Context, budgetID ID, date MonthDate) (*Month, error)
	GetLatest(ctx context.Context, budgetID ID) (*Month, error)
}
