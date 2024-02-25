package fake

import (
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"golang.org/x/net/context"
)

type categoryRepository struct{ repository }

var _ beans.CategoryRepository = (*categoryRepository)(nil)

func (r *categoryRepository) Create(ctx context.Context, tx beans.Tx, category beans.Category) error {
	r.acquire(func() { r.database.categoriesMU.RLock() })

	if _, ok := r.database.categories[category.ID]; ok {
		r.database.categoriesMU.RUnlock()
		return errors.New("duplicate")
	} else {
		r.database.categoriesMU.RUnlock()
	}

	r.txOrNow(tx, func() {
		r.acquire(func() { r.database.categoriesMU.Lock() })
		defer r.database.categoriesMU.Unlock()

		r.database.categories[category.ID] = category
	})

	return nil
}

func (r *categoryRepository) GetSingleForBudget(ctx context.Context, id beans.ID, budgetID beans.ID) (beans.Category, error) {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	defer r.database.categoriesMU.RUnlock()

	if category, ok := r.database.categories[id]; ok {
		if category.BudgetID == budgetID {
			return category, nil
		}
	}

	return beans.Category{}, beans.NewError(beans.ENOTFOUND, "category not found")
}

func (r *categoryRepository) GetCategoryGroup(ctx context.Context, id beans.ID, budgetID beans.ID) (beans.CategoryGroup, error) {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	defer r.database.categoriesMU.RUnlock()

	if group, ok := r.database.categoryGroups[id]; ok {
		if group.BudgetID == budgetID {
			return group, nil
		}
	}

	return beans.CategoryGroup{}, beans.NewError(beans.ENOTFOUND, "category group not found")
}

func (r *categoryRepository) GetCategoriesForGroup(ctx context.Context, id beans.ID, budgetID beans.ID) ([]beans.Category, error) {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	defer r.database.categoriesMU.RUnlock()

	categories := filter(values(r.database.categories), func(it beans.Category) bool { return it.BudgetID == budgetID })

	return filter(categories, func(it beans.Category) bool { return it.GroupID == id }), nil
}

func (r *categoryRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Category, error) {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	defer r.database.categoriesMU.RUnlock()

	categories := filter(values(r.database.categories), func(it beans.Category) bool { return it.BudgetID == budgetID })

	return categories, nil
}

func (r *categoryRepository) CreateGroup(ctx context.Context, tx beans.Tx, category beans.CategoryGroup) error {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	if _, ok := r.database.categoryGroups[category.ID]; ok {
		r.database.categoriesMU.RUnlock()
		return errors.New("duplicate")
	} else {
		r.database.categoriesMU.RUnlock()
	}

	r.txOrNow(tx, func() {
		r.acquire(func() { r.database.categoriesMU.Lock() })
		defer r.database.categoriesMU.Unlock()

		r.database.categoryGroups[category.ID] = category
	})

	return nil
}

func (r *categoryRepository) GetGroupsForBudget(ctx context.Context, budgetID beans.ID) ([]beans.CategoryGroup, error) {
	r.acquire(func() { r.database.categoriesMU.RLock() })
	defer r.database.categoriesMU.RUnlock()

	groups := filter(values(r.database.categoryGroups), func(it beans.CategoryGroup) bool { return it.BudgetID == budgetID })

	return groups, nil
}
