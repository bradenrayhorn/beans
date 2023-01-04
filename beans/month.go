package beans

import (
	"time"

	"golang.org/x/net/context"
)

type Month struct {
	ID       ID
	BudgetID ID
	Date     Date
}

func (m Month) String() string {
	return m.Date.Format("2006.01")
}

func NormalizeMonth(date time.Time) time.Time {
	date = date.AddDate(0, 0, -date.Day()+1)

	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

type MonthRepository interface {
	Create(ctx context.Context, tx Tx, month *Month) error
	Get(ctx context.Context, id ID) (*Month, error)
	GetByDate(ctx context.Context, budgetID ID, date time.Time) (*Month, error)
	GetLatest(ctx context.Context, budgetID ID) (*Month, error)
}

type MonthService interface {
	Get(ctx context.Context, monthID ID, budgetID ID) (*Month, error)
	GetOrCreate(ctx context.Context, budgetID ID, date time.Time) (*Month, error)
}
