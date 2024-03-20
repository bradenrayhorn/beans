package contractadapter

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/service"
	"github.com/bradenrayhorn/beans/server/specification"
)

type contractsAdapter struct {
	contracts *contract.Contracts
	services  *service.All
}

var _ specification.Interactor = (*contractsAdapter)(nil)

func New(contracts *contract.Contracts, services *service.All) specification.Interactor {
	return &contractsAdapter{contracts, services}
}

// AuthContext helpers

func (a *contractsAdapter) authContext(t *testing.T, ctx specification.Context) (*beans.AuthContext, error) {
	return a.services.User.GetAuth(context.Background(), ctx.SessionID)
}

func (a *contractsAdapter) budgetAuthContext(t *testing.T, ctx specification.Context) (*beans.BudgetAuthContext, error) {
	auth, err := a.authContext(t, ctx)
	if err != nil {
		return nil, err
	}

	budget, err := a.contracts.Budget.Get(context.Background(), auth, ctx.BudgetID)
	if err != nil {
		return nil, err
	}

	budgetAuth, err := beans.NewBudgetAuthContext(auth, budget)
	if err != nil {
		return nil, err
	}

	return budgetAuth, nil
}

// Account

func (i *contractsAdapter) AccountCreate(t *testing.T, ctx specification.Context, params beans.AccountCreate) (beans.ID, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	return i.contracts.Account.Create(context.Background(), auth, params)
}

func (i *contractsAdapter) AccountList(t *testing.T, ctx specification.Context) ([]beans.AccountWithBalance, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Account.GetAll(context.Background(), auth)
}

func (i *contractsAdapter) AccountListTransactable(t *testing.T, ctx specification.Context) ([]beans.Account, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Account.GetTransactable(context.Background(), auth)
}

func (i *contractsAdapter) AccountGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Account, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.Account{}, err
	}
	return i.contracts.Account.Get(context.Background(), auth, id)
}

// Budget

func (i *contractsAdapter) BudgetCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	auth, err := i.authContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	b, err := i.contracts.Budget.Create(context.Background(), auth, name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return b.ID, nil
}

func (i *contractsAdapter) BudgetGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Budget, error) {
	auth, err := i.authContext(t, ctx)
	if err != nil {
		return beans.Budget{}, err
	}
	return i.contracts.Budget.Get(context.Background(), auth, id)
}

func (i *contractsAdapter) BudgetGetAll(t *testing.T, ctx specification.Context) ([]beans.Budget, error) {
	auth, err := i.authContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Budget.GetAll(context.Background(), auth)
}

// Category

func (i *contractsAdapter) CategoryCreate(t *testing.T, ctx specification.Context, groupID beans.ID, name beans.Name) (beans.ID, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	category, err := i.contracts.Category.CreateCategory(context.Background(), auth, groupID, name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return category.ID, nil
}

func (i *contractsAdapter) CategoryGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Category, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.Category{}, err
	}
	return i.contracts.Category.GetCategory(context.Background(), auth, id)
}

func (i *contractsAdapter) CategoryGroupCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	group, err := i.contracts.Category.CreateGroup(context.Background(), auth, name)
	if err != nil {
		return beans.EmptyID(), err
	}
	return group.ID, nil
}

func (i *contractsAdapter) CategoryGroupGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.CategoryGroupWithCategories, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.CategoryGroupWithCategories{}, err
	}
	return i.contracts.Category.GetGroup(context.Background(), auth, id)
}

func (i *contractsAdapter) CategoryGetAll(t *testing.T, ctx specification.Context) ([]beans.CategoryGroupWithCategories, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Category.GetAll(context.Background(), auth)
}

// Month

func (i *contractsAdapter) MonthGetOrCreate(t *testing.T, ctx specification.Context, date beans.MonthDate) (beans.MonthWithDetails, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}
	return i.contracts.Month.GetOrCreate(context.Background(), auth, date)
}

func (i *contractsAdapter) MonthUpdate(t *testing.T, ctx specification.Context, id beans.ID, carryover beans.Amount) error {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return err
	}
	return i.contracts.Month.Update(context.Background(), auth, id, carryover)
}

func (i *contractsAdapter) MonthSetCategoryAmount(t *testing.T, ctx specification.Context, id beans.ID, categoryID beans.ID, amount beans.Amount) error {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return err
	}
	return i.contracts.Month.SetCategoryAmount(context.Background(), auth, id, categoryID, amount)
}

// Payee

func (i *contractsAdapter) PayeeCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	return i.contracts.Payee.CreatePayee(context.Background(), auth, name)
}

func (i *contractsAdapter) PayeeGetAll(t *testing.T, ctx specification.Context) ([]beans.Payee, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Payee.GetAll(context.Background(), auth)
}

func (i *contractsAdapter) PayeeGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Payee, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.Payee{}, err
	}
	return i.contracts.Payee.Get(context.Background(), auth, id)
}

// Transaction

func (i *contractsAdapter) TransactionCreate(t *testing.T, ctx specification.Context, params beans.TransactionCreateParams) (beans.ID, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.EmptyID(), err
	}
	return i.contracts.Transaction.Create(context.Background(), auth, params)
}

func (i *contractsAdapter) TransactionGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.TransactionWithRelations, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}
	return i.contracts.Transaction.Get(context.Background(), auth, id)
}

func (i *contractsAdapter) TransactionUpdate(t *testing.T, ctx specification.Context, params beans.TransactionUpdateParams) error {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return err
	}
	return i.contracts.Transaction.Update(context.Background(), auth, params)
}

func (i *contractsAdapter) TransactionDelete(t *testing.T, ctx specification.Context, ids []beans.ID) error {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return err
	}
	return i.contracts.Transaction.Delete(context.Background(), auth, ids)
}

func (i *contractsAdapter) TransactionGetAll(t *testing.T, ctx specification.Context) ([]beans.TransactionWithRelations, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Transaction.GetAll(context.Background(), auth)
}

func (i *contractsAdapter) TransactionGetSplits(t *testing.T, ctx specification.Context, id beans.ID) ([]beans.Split, error) {
	auth, err := i.budgetAuthContext(t, ctx)
	if err != nil {
		return nil, err
	}
	return i.contracts.Transaction.GetSplits(context.Background(), auth, id)
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

func (i *contractsAdapter) UserLogout(t *testing.T, ctx specification.Context) error {
	auth, err := i.authContext(t, ctx)
	if err != nil {
		return err
	}
	return i.contracts.User.Logout(context.Background(), auth)
}

func (i *contractsAdapter) UserGetMe(t *testing.T, ctx specification.Context) (beans.UserPublic, error) {
	auth, err := i.authContext(t, ctx)
	if err != nil {
		return beans.UserPublic{}, err
	}
	return i.contracts.User.GetMe(context.Background(), auth)
}
