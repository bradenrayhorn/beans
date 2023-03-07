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

func (d MonthDate) FirstDay() Date {
	return d.date
}

func (d MonthDate) LastDay() Date {
	return NewDate(d.Time().AddDate(0, 1, -d.Time().Day()))
}

// Creates a new MonthDate and normalizes the date.
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
	// Gets a month, its categories, and budgetable amount.
	Get(ctx context.Context, auth *BudgetAuthContext, monthID ID) (*Month, []*MonthCategory, Amount, error)

	// Creates a month in the budget with the provided date. If the month already exists, it is returned instead with no error.
	CreateMonth(ctx context.Context, auth *BudgetAuthContext, date MonthDate) (*Month, error)

	// Sets the assigned amount on a category for a month.
	SetCategoryAmount(ctx context.Context, auth *BudgetAuthContext, monthID ID, categoryID ID, amount Amount) error
}

type MonthRepository interface {
	Create(ctx context.Context, tx Tx, month *Month) error
	Get(ctx context.Context, id ID) (*Month, error)
	GetOrCreate(ctx context.Context, tx Tx, budgetID ID, date MonthDate) (*Month, error)
	GetLatest(ctx context.Context, budgetID ID) (*Month, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]*Month, error)
}
