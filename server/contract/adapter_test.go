package contract_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/stretchr/testify/require"
)

type ContractsAdapter struct {
	contracts *contract.Contracts
}

var _ specification.Interactor = (*ContractsAdapter)(nil)

// budgetAuthContext helper

func (a *ContractsAdapter) budgetAuthContext(t *testing.T, ctx specification.Context) *beans.BudgetAuthContext {
	auth, err := a.contracts.User.GetAuth(context.Background(), ctx.SessionID)
	require.Nil(t, err)

	budget, err := a.contracts.Budget.Get(context.Background(), auth, ctx.BudgetID)
	require.Nil(t, err)

	budgetAuth, err := beans.NewBudgetAuthContext(auth, budget)
	require.Nil(t, err)
	return budgetAuth
}

// userAndBudget

type userAndBudget struct {
	t         *testing.T
	sessionID beans.SessionID
	budget    beans.Budget
	context   specification.Context

	contracts  *contract.Contracts
	budgetAuth *beans.BudgetAuthContext
}

var _ specification.TestUserAndBudget = (*userAndBudget)(nil)

func (u *userAndBudget) Ctx() specification.Context {
	return u.context
}

func (u *userAndBudget) Budget() beans.Budget {
	return u.budget
}

func (u *userAndBudget) Account(opt specification.AccountOpts) beans.Account {
	name := beans.Name(beans.NewBeansID().String())
	id, err := u.contracts.Account.Create(context.Background(), u.budgetAuth, name)
	require.Nil(u.t, err)
	return beans.Account{
		ID:       id,
		Name:     name,
		BudgetID: u.budget.ID,
	}
}

func (u *userAndBudget) CategoryGroup(opt specification.CategoryGroupOpts) beans.CategoryGroup {
	name := beans.Name(beans.NewBeansID().String())
	group, err := u.contracts.Category.CreateGroup(context.Background(), u.budgetAuth, name)
	require.Nil(u.t, err)
	return group
}

func (u *userAndBudget) Category(opt specification.CategoryOpts) beans.Category {
	name := beans.Name(beans.NewBeansID().String())
	category, err := u.contracts.Category.CreateCategory(context.Background(), u.budgetAuth, opt.Group.ID, name)
	require.Nil(u.t, err)
	return category
}

func (u *userAndBudget) Transaction(opt specification.TransactionOpts) beans.Transaction {
	if opt.Date.Empty() {
		opt.Date = beans.NewDate(testutils.RandomTime())
	}
	id, err := u.contracts.Transaction.Create(context.Background(), u.budgetAuth, beans.TransactionCreateParams{
		TransactionParams: beans.TransactionParams{
			AccountID:  opt.Account.ID,
			CategoryID: opt.Category.ID,
			Amount:     opt.Amount,
			Date:       opt.Date,
		},
	})
	require.Nil(u.t, err)
	transaction, err := u.contracts.Transaction.Get(context.Background(), u.budgetAuth, id)
	require.NoError(u.t, err)

	return transaction.Transaction
}

// Test

func (i *ContractsAdapter) UserAndBudget(t *testing.T) specification.TestUserAndBudget {
	// make new user
	username := beans.NewBeansID().String()
	err := i.contracts.User.Register(
		context.Background(),
		beans.Username(username),
		beans.Password("password"),
	)
	require.Nil(t, err)

	// login as user
	session, err := i.contracts.User.Login(
		context.Background(),
		beans.Username(username),
		beans.Password("password"),
	)
	require.Nil(t, err)

	// get auth context
	auth, err := i.contracts.User.GetAuth(context.Background(), session.ID)
	require.Nil(t, err)

	// make budget
	budget, err := i.contracts.Budget.Create(context.Background(), auth, beans.Name(beans.NewBeansID().String()))
	require.NoError(t, err)

	ctx := specification.Context{SessionID: session.ID, BudgetID: budget.ID}

	return &userAndBudget{
		t:         t,
		sessionID: session.ID,
		budget:    budget,

		context:    ctx,
		contracts:  i.contracts,
		budgetAuth: i.budgetAuthContext(t, ctx),
	}
}

// Account

func (i *ContractsAdapter) AccountCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	return i.contracts.Account.Create(context.Background(), i.budgetAuthContext(t, ctx), name)
}

func (i *ContractsAdapter) AccountList(t *testing.T, ctx specification.Context) ([]beans.AccountWithBalance, error) {
	return i.contracts.Account.GetAll(context.Background(), i.budgetAuthContext(t, ctx))
}

func (i *ContractsAdapter) AccountGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Account, error) {
	return i.contracts.Account.Get(context.Background(), i.budgetAuthContext(t, ctx), id)
}
