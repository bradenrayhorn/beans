package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type monthRepository struct{ repository }

var _ beans.MonthRepository = (*monthRepository)(nil)

func (r *monthRepository) Create(ctx context.Context, tx beans.Tx, month beans.Month) error {
	r.acquire(func() { r.database.monthsMU.RLock() })
	if _, ok := r.database.months[month.ID]; ok {
		r.database.monthsMU.RUnlock()
		return errors.New("duplicate")
	} else {
		months := filter(values(r.database.months), func(it beans.Month) bool { return it.BudgetID == month.BudgetID })
		months = filter(months, func(it beans.Month) bool { return it.Date.FirstDay().Equal(month.Date.FirstDay().Time) })
		r.database.monthsMU.RUnlock()

		if len(months) > 0 {
			return errors.New("duplicate")
		}
	}

	// init carryover to 0 as postgres layer does
	month.Carryover = month.Carryover.OrZero()

	r.txOrNow(tx, func() {
		r.acquire(func() { r.database.monthsMU.Lock() })
		defer r.database.monthsMU.Unlock()

		r.database.months[month.ID] = month
	})

	return nil
}

func (r *monthRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Month, error) {
	r.acquire(func() { r.database.monthsMU.RLock() })
	defer r.database.monthsMU.RUnlock()

	if month, ok := r.database.months[id]; ok {
		if month.BudgetID == budgetID {
			return month, nil
		}
	}

	return beans.Month{}, beans.NewError(beans.ENOTFOUND, "month not found")
}

func (r *monthRepository) Update(ctx context.Context, month beans.Month) error {
	r.acquire(func() { r.database.monthsMU.Lock() })
	defer r.database.monthsMU.Unlock()

	if existing, ok := r.database.months[month.ID]; ok {
		// init carryover to zero as postgres layer does
		existing.Carryover = month.Carryover.OrZero()
		r.database.months[month.ID] = existing
		return nil
	}

	return beans.NewError(beans.ENOTFOUND, "month not found")
}

func (r *monthRepository) GetOrCreate(ctx context.Context, tx beans.Tx, budgetID beans.ID, date beans.MonthDate) (beans.Month, error) {
	r.acquire(func() { r.database.monthsMU.RLock() })
	months := filter(values(r.database.months), func(it beans.Month) bool { return it.BudgetID == budgetID })
	months = filter(months, func(it beans.Month) bool { return it.Date.FirstDay().Equal(date.FirstDay().Time) })
	r.database.monthsMU.RUnlock()

	if len(months) > 1 {
		panic("too many months")
	} else if len(months) == 1 {
		return months[0], nil
	} else {
		id := beans.NewID()
		month := beans.Month{
			ID:        id,
			BudgetID:  budgetID,
			Date:      date,
			Carryover: beans.NewAmount(0, 0),
		}

		r.txOrNow(tx, func() {
			r.acquire(func() { r.database.monthsMU.Lock() })
			defer r.database.monthsMU.Unlock()
			r.database.months[id] = month
		})

		return month, nil
	}
}

func (r *monthRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Month, error) {
	r.acquire(func() { r.database.monthsMU.RLock() })
	defer r.database.monthsMU.RUnlock()

	months := filter(values(r.database.months), func(it beans.Month) bool { return it.BudgetID == budgetID })
	return months, nil
}
