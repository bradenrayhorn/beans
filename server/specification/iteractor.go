package specification

import (
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

	// Category
	CategoryCreate(t *testing.T, ctx Context, groupID beans.ID, name beans.Name) (beans.ID, error)
	CategoryGet(t *testing.T, ctx Context, id beans.ID) (beans.Category, error)

	// CategoryGroup
	CategoryGroupCreate(t *testing.T, ctx Context, name beans.Name) (beans.ID, error)
	CategoryGroupGet(t *testing.T, ctx Context, id beans.ID) (beans.CategoryGroup, error)

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

type TransactionOpts struct {
	Account  beans.Account
	Category beans.Category
	Amount   beans.Amount
	Date     beans.Date
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
	username := beans.NewBeansID().String()
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
	budgetID, err := interactor.BudgetCreate(t, user.ctx, beans.Name(beans.NewBeansID().String()))
	require.Nil(t, err)
	budget, err := interactor.BudgetGet(t, user.ctx, budgetID)
	require.Nil(t, err)

	user.ctx.BudgetID = budget.ID

	return &userAndBudget{user, budget}
}

func (u *userAndBudget) Account(opt AccountOpts) beans.Account {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.interactor.AccountCreate(u.t, u.ctx, name)
	require.NoError(u.t, err)
	account, err := u.interactor.AccountGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return account
}

func (u *userAndBudget) CategoryGroup(opt CategoryGroupOpts) beans.CategoryGroup {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.interactor.CategoryGroupCreate(u.t, u.ctx, name)
	require.NoError(u.t, err)
	group, err := u.interactor.CategoryGroupGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return group
}

func (u *userAndBudget) Category(opt CategoryOpts) beans.Category {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.interactor.CategoryCreate(u.t, u.ctx, opt.Group.ID, name)
	require.NoError(u.t, err)
	category, err := u.interactor.CategoryGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return category
}

func (u *userAndBudget) Transaction(opt TransactionOpts) beans.Transaction {
	if opt.Date.Empty() {
		opt.Date = beans.NewDate(testutils.RandomTime())
	}
	id, err := u.interactor.TransactionCreate(u.t, u.ctx, beans.TransactionCreateParams{
		TransactionParams: beans.TransactionParams{
			AccountID:  opt.Account.ID,
			CategoryID: opt.Category.ID,
			Amount:     opt.Amount,
			Date:       opt.Date,
		},
	})
	require.NoError(u.t, err)
	transaction, err := u.interactor.TransactionGet(u.t, u.ctx, id)
	require.NoError(u.t, err)

	return transaction.Transaction
}
