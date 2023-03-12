package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func MakeUser(tb testing.TB, pool *pgxpool.Pool, username string) beans.ID {
	userID := beans.NewBeansID()
	err := postgres.NewUserRepository(pool).Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x"))
	require.Nil(tb, err)
	return userID
}

func MakeBudget(tb testing.TB, pool *pgxpool.Pool, name string, userID beans.ID) *beans.Budget {
	id := beans.NewBeansID()
	err := postgres.NewBudgetRepository(pool).Create(context.Background(), nil, id, beans.Name(name), userID)
	require.Nil(tb, err)
	return &beans.Budget{
		ID:      id,
		Name:    beans.Name(name),
		UserIDs: []beans.ID{userID},
	}
}

func MakeMonth(tb testing.TB, pool *pgxpool.Pool, budgetID beans.ID, date beans.Date) *beans.Month {
	month := &beans.Month{
		ID:        beans.NewBeansID(),
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(date),
		Carryover: beans.NewAmount(0, 0),
	}
	err := postgres.NewMonthRepository(pool).Create(context.Background(), nil, month)
	require.Nil(tb, err)
	return month
}

func MakeAccount(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) *beans.Account {
	id := beans.NewBeansID()
	err := postgres.NewAccountRepository(pool).Create(context.Background(), id, beans.Name(name), budgetID)
	require.Nil(tb, err)
	return &beans.Account{
		ID:       id,
		Name:     beans.Name(name),
		BudgetID: budgetID,
	}
}

func MakeCategoryGroup(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) *beans.CategoryGroup {
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := postgres.NewCategoryRepository(pool).CreateGroup(context.Background(), nil, group)
	require.Nil(tb, err)
	return group
}

func MakeIncomeCategoryGroup(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) *beans.CategoryGroup {
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name), IsIncome: true}
	err := postgres.NewCategoryRepository(pool).CreateGroup(context.Background(), nil, group)
	require.Nil(tb, err)
	return group
}

func MakeCategory(tb testing.TB, pool *pgxpool.Pool, name string, groupID beans.ID, budgetID beans.ID) *beans.Category {
	category := &beans.Category{ID: beans.NewBeansID(), BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)}
	err := postgres.NewCategoryRepository(pool).Create(context.Background(), nil, category)
	require.Nil(tb, err)
	return category
}

func MakeMonthCategory(tb testing.TB, pool *pgxpool.Pool, monthID beans.ID, categoryID beans.ID, amount beans.Amount) *beans.MonthCategory {
	category := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: amount}
	err := postgres.NewMonthCategoryRepository(pool).Create(context.Background(), nil, category)
	require.Nil(tb, err)
	return category
}
