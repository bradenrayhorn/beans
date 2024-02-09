package testutils

import (
	"context"
	"math/rand"
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

func (f *factory) MakeBudget(name string, userID beans.ID) beans.Budget {
	id := beans.NewBeansID()
	err := f.ds.BudgetRepository().Create(context.Background(), nil, id, beans.Name(name), userID)
	require.Nil(f.tb, err)
	return beans.Budget{
		ID:   id,
		Name: beans.Name(name),
	}
}

func (f *factory) MakeBudgetAndUser() (beans.Budget, beans.User) {
	userID := beans.NewBeansID()
	username := beans.NewBeansID().String()
	require.Nil(f.tb, f.ds.UserRepository().Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x")))

	id := beans.NewBeansID()
	budgetName := beans.NewBeansID().String()
	require.Nil(f.tb, f.ds.BudgetRepository().Create(context.Background(), nil, id, beans.Name(budgetName), userID))
	return beans.Budget{
			ID:   id,
			Name: beans.Name(budgetName),
		},
		beans.User{
			ID:           userID,
			Username:     beans.Username(username),
			PasswordHash: beans.PasswordHash("x"),
		}

}

func (f *factory) MakeMonth(budgetID beans.ID, date beans.Date) beans.Month {
	month := beans.Month{
		ID:        beans.NewBeansID(),
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(date),
		Carryover: beans.NewAmount(0, 0),
	}
	err := f.ds.MonthRepository().Create(context.Background(), nil, month)
	require.Nil(f.tb, err)
	return month
}

func (f *factory) MakeAccount(name string, budgetID beans.ID) beans.Account {
	id := beans.NewBeansID()
	err := f.ds.AccountRepository().Create(context.Background(), id, beans.Name(name), budgetID)
	require.Nil(f.tb, err)
	return beans.Account{
		ID:       id,
		Name:     beans.Name(name),
		BudgetID: budgetID,
	}
}

func (f *factory) User(user beans.User) beans.User {
	if user.ID.Empty() {
		user.ID = beans.NewBeansID()
	}

	if len(string(user.Username)) == 0 {
		user.Username = beans.Username(beans.NewBeansID().String())
	}

	if len(string(user.PasswordHash)) == 0 {
		user.PasswordHash = beans.PasswordHash("x")
	}

	require.Nil(f.tb, f.ds.UserRepository().Create(context.Background(), user.ID, user.Username, user.PasswordHash))

	return user
}

func (f *factory) Account(account beans.Account) beans.Account {
	if account.ID.Empty() {
		account.ID = beans.NewBeansID()
	}

	if len(string(account.Name)) == 0 {
		account.Name = beans.Name(beans.NewBeansID().String())
	}

	if account.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		account.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.AccountRepository().Create(context.Background(), account.ID, account.Name, account.BudgetID))

	return account
}

func (f *factory) CategoryGroup(categoryGroup beans.CategoryGroup) beans.CategoryGroup {
	if categoryGroup.ID.Empty() {
		categoryGroup.ID = beans.NewBeansID()
	}

	if len(string(categoryGroup.Name)) == 0 {
		categoryGroup.Name = beans.Name(beans.NewBeansID().String())
	}

	if categoryGroup.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		categoryGroup.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.CategoryRepository().CreateGroup(context.Background(), nil, categoryGroup))

	return categoryGroup
}

func (f *factory) Category(category beans.Category) beans.Category {
	if category.ID.Empty() {
		category.ID = beans.NewBeansID()
	}

	if len(string(category.Name)) == 0 {
		category.Name = beans.Name(beans.NewBeansID().String())
	}

	if category.GroupID.Empty() {
		categoryGroup := f.CategoryGroup(beans.CategoryGroup{BudgetID: category.BudgetID})
		category.GroupID = categoryGroup.ID
		category.BudgetID = categoryGroup.BudgetID
	}

	if category.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		category.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.CategoryRepository().Create(context.Background(), nil, category))

	return category
}

func (f *factory) MonthCategory(budgetID beans.ID, monthCategory beans.MonthCategory) beans.MonthCategory {
	if monthCategory.ID.Empty() {
		monthCategory.ID = beans.NewBeansID()
	}

	if monthCategory.MonthID.Empty() {
		month := f.Month(beans.Month{BudgetID: budgetID})
		monthCategory.MonthID = month.ID
	}

	if monthCategory.CategoryID.Empty() {
		category := f.Category(beans.Category{BudgetID: budgetID})
		monthCategory.CategoryID = category.ID
	}

	require.Nil(f.tb, f.ds.MonthCategoryRepository().Create(context.Background(), nil, monthCategory))

	return monthCategory
}

func (f *factory) Payee(payee beans.Payee) beans.Payee {
	if payee.ID.Empty() {
		payee.ID = beans.NewBeansID()
	}

	if len(string(payee.Name)) == 0 {
		payee.Name = beans.Name(beans.NewBeansID().String())
	}

	if payee.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		payee.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.PayeeRepository().Create(context.Background(), payee))

	return payee
}

func (f *factory) Transaction(budgetID beans.ID, transaction beans.Transaction) beans.Transaction {
	if transaction.ID.Empty() {
		transaction.ID = beans.NewBeansID()
	}

	if transaction.AccountID.Empty() {
		transaction.AccountID = f.Account(beans.Account{BudgetID: budgetID}).ID
	}

	if transaction.CategoryID.Empty() {
		transaction.CategoryID = f.Category(beans.Category{BudgetID: budgetID}).ID
	}

	if transaction.PayeeID.Empty() {
		transaction.PayeeID = f.Payee(beans.Payee{BudgetID: budgetID}).ID
	}

	if transaction.Date.Empty() {
		transaction.Date = beans.NewDate(RandomTime())
	}

	if transaction.Amount.Empty() {
		coefficient := rand.Int63n(10)
		exponent := rand.Int31n(2)
		if coefficient == 0 {
			exponent = 0
		}
		transaction.Amount = beans.NewAmount(coefficient, exponent)
	}

	require.Nil(f.tb, f.ds.TransactionRepository().Create(context.Background(), transaction))

	return transaction
}

func (f *factory) Month(month beans.Month) beans.Month {
	if month.ID.Empty() {
		month.ID = beans.NewBeansID()
	}

	if month.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		month.BudgetID = defaultBudget.ID
	}

	if month.Date.Empty() {
		month.Date = beans.NewMonthDate(beans.NewDate(RandomTime()))
	}

	if month.Carryover.Empty() {
		month.Carryover = beans.NewAmount(rand.Int63n(10), rand.Int31n(2))
	}

	require.Nil(f.tb, f.ds.MonthRepository().Create(context.Background(), nil, month))

	return month
}

func (f *factory) MakeCategoryGroup(name string, budgetID beans.ID) beans.CategoryGroup {
	group := beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *factory) MakeIncomeCategoryGroup(name string, budgetID beans.ID) beans.CategoryGroup {
	group := beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name), IsIncome: true}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *factory) MakeCategory(name string, groupID beans.ID, budgetID beans.ID) beans.Category {
	category := beans.Category{ID: beans.NewBeansID(), BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *factory) MakeMonthCategory(monthID beans.ID, categoryID beans.ID, amount beans.Amount) beans.MonthCategory {
	category := beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: amount}
	err := f.ds.MonthCategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *factory) MakePayee(name string, budgetID beans.ID) beans.Payee {
	payee := beans.Payee{ID: beans.NewBeansID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.PayeeRepository().Create(context.Background(), payee)
	require.Nil(f.tb, err)
	return payee
}
