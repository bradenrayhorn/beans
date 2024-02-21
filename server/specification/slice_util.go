package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
)

func mustFind[T any](t *testing.T, items []T, matches func(a T) bool) T {
	for _, item := range items {
		if matches(item) {
			return item
		}
	}

	t.Error("Could not find item in list.")

	var empty T
	return empty
}

func findAccount(t *testing.T, items []beans.AccountWithBalance, id beans.ID, do func(account beans.AccountWithBalance)) {
	account := mustFind(t, items, func(a beans.AccountWithBalance) bool { return a.ID == id })
	do(account)
}

func findCategoryGroup(t *testing.T, items []beans.CategoryGroupWithCategories, id beans.ID, do func(it beans.CategoryGroupWithCategories)) {
	res := mustFind(t, items, func(a beans.CategoryGroupWithCategories) bool { return a.ID == id })
	do(res)
}

func findMonthCategory(t *testing.T, items []beans.MonthCategoryWithDetails, categoryID beans.ID, do func(it beans.MonthCategoryWithDetails)) {
	res := mustFind(t, items, func(a beans.MonthCategoryWithDetails) bool { return a.CategoryID == categoryID })
	do(res)
}
