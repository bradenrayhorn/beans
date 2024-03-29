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

func findAccountWithBalance(t *testing.T, items []beans.AccountWithBalance, id beans.ID, do func(it beans.AccountWithBalance)) {
	res := mustFind(t, items, func(a beans.AccountWithBalance) bool { return a.ID == id })
	do(res)
}

func findAccount(t *testing.T, items []beans.Account, id beans.ID, do func(it beans.Account)) {
	res := mustFind(t, items, func(a beans.Account) bool { return a.ID == id })
	do(res)
}

func findCategoryGroup(t *testing.T, items []beans.CategoryGroupWithCategories, id beans.ID, do func(it beans.CategoryGroupWithCategories)) {
	res := mustFind(t, items, func(a beans.CategoryGroupWithCategories) bool { return a.ID == id })
	do(res)
}

func findMonthCategory(t *testing.T, items []beans.MonthCategoryWithDetails, categoryID beans.ID, do func(it beans.MonthCategoryWithDetails)) {
	res := mustFind(t, items, func(a beans.MonthCategoryWithDetails) bool { return a.CategoryID == categoryID })
	do(res)
}

func findPayee(t *testing.T, items []beans.Payee, id beans.ID, do func(it beans.Payee)) {
	res := mustFind(t, items, func(a beans.Payee) bool { return a.ID == id })
	do(res)
}

func findSplit(t *testing.T, items []beans.Split, id beans.ID, do func(it beans.Split)) {
	res := mustFind(t, items, func(a beans.Split) bool { return a.ID == id })
	do(res)
}

func findTransaction(t *testing.T, items []beans.TransactionWithRelations, id beans.ID, do func(it beans.TransactionWithRelations)) {
	res := mustFind(t, items, func(a beans.TransactionWithRelations) bool { return a.ID == id })
	do(res)
}
