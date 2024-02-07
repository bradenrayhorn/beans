package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustFind[T any](t *testing.T, items []T, matches func(a T) bool) T {
	for _, item := range items {
		if matches(item) {
			return item
		}
	}

	t.Error("Could not find item in list.")

	var empty T
	return empty
}

func findAccount(t *testing.T, items []beans.Account, id beans.ID, do func(account beans.Account)) {
	account := mustFind(t, items, func(a beans.Account) bool { return a.ID == id })
	do(account)
}

func TestAccounts(t *testing.T, interactor Interactor) {

	t.Run("create", func(t *testing.T) {
		t.Run("cannot create with invalid name", func(t *testing.T) {
			c := interactor.UserAndBudget(t)

			_, err := interactor.AccountCreate(t, c.Ctx(), "")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := interactor.UserAndBudget(t)

			// create account
			accountID, err := interactor.AccountCreate(t, c.Ctx(), "New Account")
			require.NoError(t, err)

			// check if account was saved properly
			account, err := interactor.AccountGet(t, c.Ctx(), accountID)
			require.NoError(t, err)

			assert.Equal(t, beans.Name("New Account"), account.Name)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get a non-existent account", func(t *testing.T) {
			c := interactor.UserAndBudget(t)

			_, err := interactor.AccountGet(t, c.Ctx(), beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get account from another budget", func(t *testing.T) {
			c1 := interactor.UserAndBudget(t)
			c2 := interactor.UserAndBudget(t)

			// Try and get an account from budget 2 with budget 1
			account := c2.Account(AccountOpts{})

			_, err := interactor.AccountGet(t, c1.Ctx(), account.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get all", func(t *testing.T) {
		c := interactor.UserAndBudget(t)

		account1 := c.Account(AccountOpts{})
		account2 := c.Account(AccountOpts{})

		// add a +$6 transaction to account 1
		categoryGroup := c.CategoryGroup(CategoryGroupOpts{})
		category := c.Category(CategoryOpts{Group: categoryGroup})
		c.Transaction(TransactionOpts{
			Account:  account1,
			Category: category,
			Amount:   beans.NewAmount(6, 0),
		})

		// list accounts, check if accounts are proper
		accounts, err := interactor.AccountList(t, c.Ctx())
		require.NoError(t, err)
		require.Len(t, accounts, 2)

		findAccount(t, accounts, account1.ID, func(account beans.Account) {
			assert.Equal(t, account1.Name, account.Name)
			assert.Equal(t, beans.NewAmount(6, 0), account.Balance)
		})
		findAccount(t, accounts, account2.ID, func(account beans.Account) {
			assert.Equal(t, account2.Name, account.Name)
			assert.Equal(t, beans.NewAmount(0, 0), account.Balance)
		})
	})
}
