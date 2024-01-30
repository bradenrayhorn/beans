package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/require"
)

type factory struct {
	tb testing.TB
	ds beans.DataSource
}

func Factory(tb testing.TB, ds beans.DataSource) *factory {
	return &factory{tb, ds}
}

func (f *factory) MakeUser(username string) beans.ID {
	userID := beans.NewBeansID()
	err := f.ds.UserRepository().Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x"))
	require.Nil(f.tb, err)
	return userID
}

func (f *factory) MakeBudget(name string, userID beans.ID) *beans.Budget {
	id := beans.NewBeansID()
	err := f.ds.BudgetRepository().Create(context.Background(), nil, id, beans.Name(name), userID)
	require.Nil(f.tb, err)
	return &beans.Budget{
		ID:      id,
		Name:    beans.Name(name),
		UserIDs: []beans.ID{userID},
	}
}

func (f *factory) MakeMonth(budgetID beans.ID, date beans.Date) *beans.Month {
	month := &beans.Month{
		ID:        beans.NewBeansID(),
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(date),
		Carryover: beans.NewAmount(0, 0),
	}
	err := f.ds.MonthRepository().Create(context.Background(), nil, month)
	require.Nil(f.tb, err)
	return month
}

func (f *factory) MakeAccount(name string, budgetID beans.ID) *beans.Account {
	id := beans.NewBeansID()
	err := f.ds.AccountRepository().Create(context.Background(), id, beans.Name(name), budgetID)
	require.Nil(f.tb, err)
	return &beans.Account{
		ID:       id,
		Name:     beans.Name(name),
		BudgetID: budgetID,
	}
}

func (f *factory) MakeCategoryGroup(name string, budgetID beans.ID) *beans.CategoryGroup {
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *factory) MakeIncomeCategoryGroup(name string, budgetID beans.ID) *beans.CategoryGroup {
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name), IsIncome: true}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *factory) MakeCategory(name string, groupID beans.ID, budgetID beans.ID) *beans.Category {
	category := &beans.Category{ID: beans.NewBeansID(), BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *factory) MakeMonthCategory(monthID beans.ID, categoryID beans.ID, amount beans.Amount) *beans.MonthCategory {
	category := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: amount}
	err := f.ds.MonthCategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *factory) MakePayee(name string, budgetID beans.ID) *beans.Payee {
	payee := &beans.Payee{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.PayeeRepository().Create(context.Background(), payee)
	require.Nil(f.tb, err)
	return payee
}
