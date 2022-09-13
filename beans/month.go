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

type MonthRepository interface {
	Create(ctx context.Context, month *Month) error
	GetByDate(ctx context.Context, budgetID ID, date time.Time) (*Month, error)
}

type MonthService interface {
	GetOrCreate(ctx context.Context, budgetID ID, date time.Time) (*Month, error)
}
