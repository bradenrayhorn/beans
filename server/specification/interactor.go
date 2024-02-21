package specification

import (
	"encoding/json"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/require"
)

type Interactor interface {
	// Account
	AccountCreate(t *testing.T, ctx Context, name beans.Name) (beans.ID, error)
	AccountList(t *testing.T, ctx Context) ([]beans.AccountWithBalance, error)
	AccountGet(t *testing.T, ctx Context, id beans.ID) (beans.Account, error)

	// Budget
	BudgetCreate(t *testing.T, ctx Context, name beans.Name) (beans.ID, error)
	BudgetGet(t *testing.T, ctx Context, id beans.ID) (beans.Budget, error)
	BudgetGetAll(t *testing.T, ctx Context) ([]beans.Budget, error)

	// Category
	CategoryCreate(t *testing.T, ctx Context, groupID beans.ID, name beans.Name) (beans.ID, error)
	CategoryGet(t *testing.T, ctx Context, id beans.ID) (beans.Category, error)

	CategoryGroupCreate(t *testing.T, ctx Context, name beans.Name) (beans.ID, error)
	CategoryGroupGet(t *testing.T, ctx Context, id beans.ID) (beans.CategoryGroupWithCategories, error)

	CategoryGetAll(t *testing.T, ctx Context) ([]beans.CategoryGroupWithCategories, error)

	// Month
	MonthGetOrCreate(t *testing.T, ctx Context, date beans.MonthDate) (beans.MonthWithDetails, error)
	MonthUpdate(t *testing.T, ctx Context, monthID beans.ID, carryover beans.Amount) error
	MonthSetCategoryAmount(t *testing.T, ctx Context, monthID beans.ID, categoryID beans.ID, amount beans.Amount) error

	// Transaction
	TransactionCreate(t *testing.T, ctx Context, params beans.TransactionCreateParams) (beans.ID, error)
	TransactionGet(t *testing.T, ctx Context, id beans.ID) (beans.TransactionWithRelations, error)

	// User
	UserRegister(t *testing.T, ctx Context, username beans.Username, password beans.Password) error
	UserLogin(t *testing.T, ctx Context, username beans.Username, password beans.Password) (beans.SessionID, error)
}

// Common parameters that need to be passed on most requests.
type Context struct {
	SessionID beans.SessionID
	BudgetID  beans.ID
}

type AccountOpts struct {
}

type CategoryGroupOpts struct {
}

type CategoryOpts struct {
	Group beans.CategoryGroup
}

type MonthOpts struct {
	Date      string
	Carryover string
}

type TransactionOpts struct {
	Account  beans.Account
	Category beans.Category
	Amount   string
	Date     string
}

type user struct {
	t         *testing.T
	sessionID beans.SessionID
	ctx       Context

	interactor Interactor
}

type userAndBudget struct {
	*user
	budget beans.Budget
}

func makeUser(t *testing.T, interactor Interactor) *user {
	ctx := Context{}

	// make new user
	username := beans.NewID().String()
	err := interactor.UserRegister(t, ctx, beans.Username(username), beans.Password("password"))
	require.Nil(t, err)

	// login as user
	sessionID, err := interactor.UserLogin(t, ctx, beans.Username(username), beans.Password("password"))
	require.Nil(t, err)

	ctx.SessionID = sessionID

	return &user{t, sessionID, ctx, interactor}
}

func makeUserAndBudget(t *testing.T, interactor Interactor) *userAndBudget {
	user := makeUser(t, interactor)

	// make budget
	budgetID, err := interactor.BudgetCreate(t, user.ctx, beans.Name(beans.NewID().String()))
	require.Nil(t, err)
	budget, err := interactor.BudgetGet(t, user.ctx, budgetID)
	require.Nil(t, err)

	user.ctx.BudgetID = budget.ID

	return &userAndBudget{user, budget}
}

// Factory functions

func (u *userAndBudget) Account(opt AccountOpts) beans.Account {
	name := beans.Name(beans.NewID().String())

	id, err := u.interactor.AccountCreate(u.t, u.ctx, name)
	require.NoError(u.t, err)
	account, err := u.interactor.AccountGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return account
}

func (u *userAndBudget) CategoryGroup(opt CategoryGroupOpts) beans.CategoryGroup {
	name := beans.Name(beans.NewID().String())

	id, err := u.interactor.CategoryGroupCreate(u.t, u.ctx, name)
	require.NoError(u.t, err)
	group, err := u.interactor.CategoryGroupGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return group.CategoryGroup
}

func (u *userAndBudget) Category(opt CategoryOpts) beans.Category {
	name := beans.Name(beans.NewID().String())

	if opt.Group.ID.Empty() {
		opt.Group = u.CategoryGroup(CategoryGroupOpts{})
	}

	id, err := u.interactor.CategoryCreate(u.t, u.ctx, opt.Group.ID, name)
	require.NoError(u.t, err)
	category, err := u.interactor.CategoryGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return category
}

func (u *userAndBudget) Month(opt MonthOpts) beans.Month {
	date := testutils.NewMonthDate(u.t, opt.Date)

	month, err := u.interactor.MonthGetOrCreate(u.t, u.ctx, date)
	require.NoError(u.t, err)

	if opt.Carryover != "" {
		var carryover beans.Amount
		require.NoError(u.t, json.Unmarshal([]byte(opt.Carryover), &carryover))

		err = u.interactor.MonthUpdate(u.t, u.ctx, month.ID, carryover)
		require.NoError(u.t, err)
	}

	month, err = u.interactor.MonthGetOrCreate(u.t, u.ctx, date)
	require.NoError(u.t, err)

	return month.Month
}

func (u *userAndBudget) Transaction(opt TransactionOpts) beans.Transaction {
	params := beans.TransactionParams{
		AccountID:  opt.Account.ID,
		CategoryID: opt.Category.ID,
	}

	// date
	if opt.Date == "" {
		params.Date = beans.NewDate(testutils.RandomTime())
	} else {
		params.Date = testutils.NewDate(u.t, opt.Date)
	}

	// amount
	require.NoError(u.t, json.Unmarshal([]byte(opt.Amount), &params.Amount))

	// create
	id, err := u.interactor.TransactionCreate(u.t, u.ctx, beans.TransactionCreateParams{
		TransactionParams: params,
	})
	require.NoError(u.t, err)
	transaction, err := u.interactor.TransactionGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return transaction.Transaction
}

// Other helpers

func (u *userAndBudget) findIncomeCategory() beans.Category {
	groups, err := u.interactor.CategoryGetAll(u.t, u.ctx)
	require.NoError(u.t, err)

	for _, group := range groups {
		if group.IsIncome {
			category, err := u.interactor.CategoryGet(u.t, u.ctx, group.Categories[0].ID)
			require.NoError(u.t, err)
			return category
		}
	}

	u.t.Fatal("could not find income category")

	return beans.Category{}
}

func (u *userAndBudget) setAssigned(month beans.Month, category beans.Category, amount string) {
	// amount
	var amountObj beans.Amount
	require.NoError(u.t, json.Unmarshal([]byte(amount), &amountObj))

	err := u.interactor.MonthSetCategoryAmount(u.t, u.ctx, month.ID, category.ID, amountObj)
	require.NoError(u.t, err)
}
