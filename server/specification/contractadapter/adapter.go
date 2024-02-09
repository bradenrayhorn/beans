package contractadapter

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/stretchr/testify/require"
)

type contractsAdapter struct {
	contracts *contract.Contracts
}

var _ specification.Interactor = (*contractsAdapter)(nil)

func New(contracts *contract.Contracts) specification.Interactor {
	return &contractsAdapter{contracts}
}

// AuthContext helpers

func (a *contractsAdapter) authContext(t *testing.T, ctx specification.Context) *beans.AuthContext {
	auth, err := a.contracts.User.GetAuth(context.Background(), ctx.SessionID)
	require.NoError(t, err)
	return auth
}

func (a *contractsAdapter) budgetAuthContext(t *testing.T, ctx specification.Context) *beans.BudgetAuthContext {
	auth := a.authContext(t, ctx)

	budget, err := a.contracts.Budget.Get(context.Background(), auth, ctx.BudgetID)
	require.Nil(t, err)

	budgetAuth, err := beans.NewBudgetAuthContext(auth, budget)
	require.Nil(t, err)
	return budgetAuth
}

// Account

func (i *contractsAdapter) AccountCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	return i.contracts.Account.Create(context.Background(), i.budgetAuthContext(t, ctx), name)
}

func (i *contractsAdapter) AccountList(t *testing.T, ctx specification.Context) ([]beans.AccountWithBalance, error) {
	return i.contracts.Account.GetAll(context.Background(), i.budgetAuthContext(t, ctx))
}

func (i *contractsAdapter) AccountGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Account, error) {
	return i.contracts.Account.Get(context.Background(), i.budgetAuthContext(t, ctx), id)
}

// Budget

func (i *contractsAdapter) BudgetCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	b, err := i.contracts.Budget.Create(context.Background(), i.authContext(t, ctx), name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return b.ID, nil
}

func (i *contractsAdapter) BudgetGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Budget, error) {
	return i.contracts.Budget.Get(context.Background(), i.authContext(t, ctx), id)
}

// CategoryGroup

func (i *contractsAdapter) CategoryGroupCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	group, err := i.contracts.Category.CreateGroup(context.Background(), i.budgetAuthContext(t, ctx), name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return group.ID, nil
}

func (i *contractsAdapter) CategoryGroupGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.CategoryGroup, error) {
	group, err := i.contracts.Category.GetGroup(context.Background(), i.budgetAuthContext(t, ctx), id)
	if err != nil {
		return beans.CategoryGroup{}, err
	}
	return group.CategoryGroup, nil
}

// Category

func (i *contractsAdapter) CategoryCreate(t *testing.T, ctx specification.Context, groupID beans.ID, name beans.Name) (beans.ID, error) {
	category, err := i.contracts.Category.CreateCategory(context.Background(), i.budgetAuthContext(t, ctx), groupID, name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return category.ID, nil
}

func (i *contractsAdapter) CategoryGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Category, error) {
	return i.contracts.Category.GetCategory(context.Background(), i.budgetAuthContext(t, ctx), id)
}

// Transaction

func (i *contractsAdapter) TransactionCreate(t *testing.T, ctx specification.Context, params beans.TransactionCreateParams) (beans.ID, error) {
	return i.contracts.Transaction.Create(context.Background(), i.budgetAuthContext(t, ctx), params)
}

func (i *contractsAdapter) TransactionGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.TransactionWithRelations, error) {
	return i.contracts.Transaction.Get(context.Background(), i.budgetAuthContext(t, ctx), id)
}

// User

func (i *contractsAdapter) UserRegister(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) error {
	return i.contracts.User.Register(context.Background(), username, password)
}

func (i *contractsAdapter) UserLogin(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) (beans.SessionID, error) {
	session, err := i.contracts.User.Login(context.Background(), username, password)
	if err != nil {
		return "", err
	}
	return session.ID, nil
}
