package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
)

type Interactor interface {
	// Test
	UserAndBudget(t *testing.T) TestUserAndBudget

	// Account
	AccountCreate(t *testing.T, ctx Context, name beans.Name) (beans.ID, error)
	AccountList(t *testing.T, ctx Context) ([]beans.AccountWithBalance, error)
	AccountGet(t *testing.T, ctx Context, id beans.ID) (beans.Account, error)
}

// Common parameters that need to be passed on most requests.
type Context struct {
	SessionID beans.SessionID
	BudgetID  beans.ID
}

type TestUserAndBudget interface {
	Ctx() Context
	Budget() beans.Budget

	Account(opt AccountOpts) beans.Account
	CategoryGroup(opt CategoryGroupOpts) beans.CategoryGroup
	Category(opt CategoryOpts) beans.Category
	Transaction(opt TransactionOpts) beans.Transaction
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
