package beans

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
)

type Month struct {
	ID        ID
	BudgetID  ID
	Date      MonthDate
	Carryover Amount
}

type MonthWithDetails struct {
	Month
	CarriedOver Amount
	Income      Amount
	Assigned    Amount

	Budgetable Amount
	Categories []MonthCategoryWithDetails
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

func (d MonthDate) Previous() MonthDate {
	return NewMonthDate(NewDate(d.FirstDay().AddDate(0, -1, 0)))
}

func (d MonthDate) Empty() bool {
	return d.date.Empty()
}

func (d *MonthDate) UnmarshalJSON(b []byte) error {
	var date Date
	if err := json.Unmarshal(b, &date); err != nil {
		return err
	}

	*d = NewMonthDate(date)
	return nil
}

func (d MonthDate) MarshalJSON() ([]byte, error) {
	return d.date.MarshalJSON()
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
	// If the month does not exist it is created.
	GetOrCreate(ctx context.Context, auth *BudgetAuthContext, date MonthDate) (MonthWithDetails, error)

	// Updates the given month.
	Update(ctx context.Context, auth *BudgetAuthContext, monthID ID, carryover Amount) error

	// Sets the assigned amount on a category for a month.
	SetCategoryAmount(ctx context.Context, auth *BudgetAuthContext, monthID ID, categoryID ID, amount Amount) error
}

type MonthRepository interface {
	Create(ctx context.Context, tx Tx, month Month) error
	Get(ctx context.Context, budgetID ID, id ID) (Month, error)
	// Only updates the Carryover field.
	Update(ctx context.Context, month Month) error
	GetOrCreate(ctx context.Context, tx Tx, budgetID ID, date MonthDate) (Month, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]Month, error)
}
