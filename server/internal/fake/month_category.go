package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type monthCategoryRepository struct{ repository }

var _ beans.MonthCategoryRepository = (*monthCategoryRepository)(nil)

func (r *monthCategoryRepository) Create(ctx context.Context, tx beans.Tx, monthCategory beans.MonthCategory) error {
	r.acquire(func() { r.database.monthCategoriesMU.RLock() })

	existing := filter(values(r.database.monthCategories), func(it beans.MonthCategory) bool {
		return (it.MonthID == monthCategory.MonthID && it.CategoryID == monthCategory.CategoryID) || it.ID == monthCategory.ID
	})

	if len(existing) > 0 {
		r.database.monthCategoriesMU.RUnlock()
		return errors.New("duplicate")
	} else {
		r.database.monthCategoriesMU.RUnlock()
	}

	r.txOrNow(tx, func() {
		r.acquire(func() { r.database.monthCategoriesMU.Lock() })
		defer r.database.monthCategoriesMU.Unlock()

		// init amount to zero if it is not set
		// this is what the postgres mapper does
		monthCategory.Amount = monthCategory.Amount.OrZero()

		r.database.monthCategories[monthCategory.ID] = monthCategory
	})

	return nil
}

func (r *monthCategoryRepository) UpdateAmount(ctx context.Context, monthCategory beans.MonthCategory) error {
	r.acquire(func() { r.database.monthCategoriesMU.Lock() })
	defer r.database.monthCategoriesMU.Unlock()

	if existing, ok := r.database.monthCategories[monthCategory.ID]; ok {
		existing.Amount = monthCategory.Amount
		r.database.monthCategories[monthCategory.ID] = existing
		return nil
	}

	return beans.NewError(beans.ENOTFOUND, "month category not found")
}

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, month beans.Month) ([]beans.MonthCategory, error) {
	r.acquire(func() { r.database.monthCategoriesMU.RLock() })
	defer r.database.monthCategoriesMU.RUnlock()

	return filter(values(r.database.monthCategories), func(it beans.MonthCategory) bool {
		return it.MonthID == month.ID
	}), nil
}

func (r *monthCategoryRepository) GetAssignedByCategory(ctx context.Context, budgetID beans.ID, before beans.Date) (map[beans.ID]beans.Amount, error) {
	r.acquire(func() {
		r.database.monthCategoriesMU.RLock()
		r.database.monthsMU.RLock()
	})
	defer r.database.monthCategoriesMU.RUnlock()
	defer r.database.monthsMU.RUnlock()

	categories := filter(values(r.database.monthCategories), func(it beans.MonthCategory) bool {
		if m, ok := r.database.months[it.MonthID]; ok {
			if m.BudgetID == budgetID {
				return m.Date.FirstDay().Before(before.Time)
			}
		}
		return false
	})

	assignedByCategory := make(map[beans.ID]beans.Amount)
	for _, c := range categories {
		if current, ok := assignedByCategory[c.CategoryID]; ok {
			sum, err := beans.Arithmetic.Add(current, c.Amount)
			if err != nil {
				panic(err)
			}
			assignedByCategory[c.CategoryID] = sum
		} else {
			assignedByCategory[c.CategoryID] = c.Amount
		}
	}

	return assignedByCategory, nil
}

func (r *monthCategoryRepository) GetOrCreate(ctx context.Context, tx beans.Tx, month beans.Month, categoryID beans.ID) (beans.MonthCategory, error) {
	r.acquire(func() { r.database.monthCategoriesMU.RLock() })
	categories := filter(values(r.database.monthCategories), func(it beans.MonthCategory) bool { return it.MonthID == month.ID && it.CategoryID == categoryID })
	r.database.monthCategoriesMU.RUnlock()

	if len(categories) > 1 {
		panic("too many categories")
	} else if len(categories) == 1 {
		return categories[0], nil
	} else {
		id := beans.NewID()
		res := beans.MonthCategory{
			ID:         id,
			MonthID:    month.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(0, 0),
		}

		r.txOrNow(tx, func() {
			r.acquire(func() { r.database.monthCategoriesMU.Lock() })
			defer r.database.monthCategoriesMU.Unlock()

			r.database.monthCategories[id] = res
		})

		return res, nil
	}
}

func (r *monthCategoryRepository) GetAssignedInMonth(ctx context.Context, month beans.Month) (beans.Amount, error) {
	r.acquire(func() { r.database.monthCategoriesMU.RLock() })
	defer r.database.monthCategoriesMU.RUnlock()

	categories := filter(values(r.database.monthCategories), func(it beans.MonthCategory) bool { return it.MonthID == month.ID })

	return reduce(categories, beans.NewAmount(0, 0), func(it beans.MonthCategory, acc beans.Amount) beans.Amount {
		r, err := beans.Arithmetic.Add(acc, it.Amount)
		if err != nil {
			panic(err)
		}
		return r
	}), nil
}
