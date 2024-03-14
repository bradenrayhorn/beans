package testutils

import (
	"context"
	"math/rand"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/require"
)

type Factory struct {
	tb testing.TB
	ds beans.DataSource
}

func NewFactory(tb testing.TB, ds beans.DataSource) *Factory {
	return &Factory{tb, ds}
}

func (f *Factory) MakeUser(username string) beans.ID {
	userID := beans.NewID()
	err := f.ds.UserRepository().Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x"))
	require.Nil(f.tb, err)
	return userID
}

func (f *Factory) MakeBudget(name string, userID beans.ID) beans.Budget {
	id := beans.NewID()
	err := f.ds.BudgetRepository().Create(context.Background(), nil, id, beans.Name(name), userID)
	require.Nil(f.tb, err)
	return beans.Budget{
		ID:   id,
		Name: beans.Name(name),
	}
}

func (f *Factory) MakeBudgetAndUser() (beans.Budget, beans.User) {
	userID := beans.NewID()
	username := beans.NewID().String()
	require.Nil(f.tb, f.ds.UserRepository().Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x")))

	id := beans.NewID()
	budgetName := beans.NewID().String()
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

func (f *Factory) MakeMonth(budgetID beans.ID, date beans.Date) beans.Month {
	month := beans.Month{
		ID:        beans.NewID(),
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(date),
		Carryover: beans.NewAmount(0, 0),
	}
	err := f.ds.MonthRepository().Create(context.Background(), nil, month)
	require.Nil(f.tb, err)
	return month
}

func (f *Factory) User(user beans.User) beans.User {
	if user.ID.Empty() {
		user.ID = beans.NewID()
	}

	if len(string(user.Username)) == 0 {
		user.Username = beans.Username(beans.NewID().String())
	}

	if len(string(user.PasswordHash)) == 0 {
		user.PasswordHash = beans.PasswordHash("x")
	}

	require.Nil(f.tb, f.ds.UserRepository().Create(context.Background(), user.ID, user.Username, user.PasswordHash))

	return user
}

func (f *Factory) Account(account beans.Account) beans.Account {
	if account.ID.Empty() {
		account.ID = beans.NewID()
	}

	if len(string(account.Name)) == 0 {
		account.Name = beans.Name(beans.NewID().String())
	}

	if account.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		account.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.AccountRepository().Create(context.Background(), account))

	return account
}

func (f *Factory) CategoryGroup(categoryGroup beans.CategoryGroup) beans.CategoryGroup {
	if categoryGroup.ID.Empty() {
		categoryGroup.ID = beans.NewID()
	}

	if len(string(categoryGroup.Name)) == 0 {
		categoryGroup.Name = beans.Name(beans.NewID().String())
	}

	if categoryGroup.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		categoryGroup.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.CategoryRepository().CreateGroup(context.Background(), nil, categoryGroup))

	return categoryGroup
}

func (f *Factory) Category(category beans.Category) beans.Category {
	if category.ID.Empty() {
		category.ID = beans.NewID()
	}

	if len(string(category.Name)) == 0 {
		category.Name = beans.Name(beans.NewID().String())
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

func (f *Factory) MonthCategory(budgetID beans.ID, monthCategory beans.MonthCategory) beans.MonthCategory {
	if monthCategory.ID.Empty() {
		monthCategory.ID = beans.NewID()
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

func (f *Factory) Payee(payee beans.Payee) beans.Payee {
	if payee.ID.Empty() {
		payee.ID = beans.NewID()
	}

	if len(string(payee.Name)) == 0 {
		payee.Name = beans.Name(beans.NewID().String())
	}

	if payee.BudgetID.Empty() {
		defaultBudget, _ := f.MakeBudgetAndUser()
		payee.BudgetID = defaultBudget.ID
	}

	require.Nil(f.tb, f.ds.PayeeRepository().Create(context.Background(), payee))

	return payee
}

func (f *Factory) Transaction(budgetID beans.ID, transaction beans.Transaction) beans.Transaction {
	if transaction.ID.Empty() {
		transaction.ID = beans.NewID()
	}

	if transaction.AccountID.Empty() {
		transaction.AccountID = f.Account(beans.Account{BudgetID: budgetID}).ID
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

	require.Nil(f.tb, f.ds.TransactionRepository().Create(
		context.Background(),
		[]beans.Transaction{transaction},
	))

	return transaction
}

func (f *Factory) Transfer(budgetID beans.ID, accountA beans.Account, accountB beans.Account, amount beans.Amount) []beans.Transaction {
	date := beans.NewDate(RandomTime())

	transactionA := beans.Transaction{
		ID:        beans.NewID(),
		AccountID: accountA.ID,
		Amount:    amount,
		Date:      date,
	}
	transactionB := beans.Transaction{
		ID:        beans.NewID(),
		AccountID: accountB.ID,
		Amount:    beans.Arithmetic.Negate(amount),
		Date:      date,
	}

	transactionA.TransferID = transactionB.ID
	transactionB.TransferID = transactionA.ID

	require.Nil(f.tb, f.ds.TransactionRepository().Create(
		context.Background(),
		[]beans.Transaction{transactionA, transactionB},
	))

	return []beans.Transaction{transactionA, transactionB}
}

func (f *Factory) Month(month beans.Month) beans.Month {
	if month.ID.Empty() {
		month.ID = beans.NewID()
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

func (f *Factory) MakeCategoryGroup(name string, budgetID beans.ID) beans.CategoryGroup {
	group := beans.CategoryGroup{ID: beans.NewID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *Factory) MakeIncomeCategoryGroup(name string, budgetID beans.ID) beans.CategoryGroup {
	group := beans.CategoryGroup{ID: beans.NewID(), BudgetID: budgetID, Name: beans.Name(name), IsIncome: true}
	err := f.ds.CategoryRepository().CreateGroup(context.Background(), nil, group)
	require.Nil(f.tb, err)
	return group
}

func (f *Factory) MakeCategory(name string, groupID beans.ID, budgetID beans.ID) beans.Category {
	category := beans.Category{ID: beans.NewID(), BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)}
	err := f.ds.CategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *Factory) MakeMonthCategory(monthID beans.ID, categoryID beans.ID, amount beans.Amount) beans.MonthCategory {
	category := beans.MonthCategory{ID: beans.NewID(), MonthID: monthID, CategoryID: categoryID, Amount: amount}
	err := f.ds.MonthCategoryRepository().Create(context.Background(), nil, category)
	require.Nil(f.tb, err)
	return category
}

func (f *Factory) MakePayee(name string, budgetID beans.ID) beans.Payee {
	payee := beans.Payee{ID: beans.NewID(), BudgetID: budgetID, Name: beans.Name(name)}
	err := f.ds.PayeeRepository().Create(context.Background(), payee)
	require.Nil(f.tb, err)
	return payee
}
