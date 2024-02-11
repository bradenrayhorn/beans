package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAccount(t *testing.T, interactor Interactor) {

	t.Run("create", func(t *testing.T) {
		t.Run("cannot create with invalid name", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.AccountCreate(t, c.ctx, "")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create account
			accountID, err := interactor.AccountCreate(t, c.ctx, "New Account")
			require.NoError(t, err)

			// check if account was saved properly
			account, err := interactor.AccountGet(t, c.ctx, accountID)
			require.NoError(t, err)

			assert.False(t, account.ID.Empty())
			assert.Equal(t, beans.Name("New Account"), account.Name)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get a non-existent account", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.AccountGet(t, c.ctx, beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get account from another budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			// Try and get an account from budget 2 with budget 1
			account := c2.Account(AccountOpts{})

			_, err := interactor.AccountGet(t, c1.ctx, account.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get all", func(t *testing.T) {
		c := makeUserAndBudget(t, interactor)

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
		accounts, err := interactor.AccountList(t, c.ctx)
		require.NoError(t, err)
		require.Len(t, accounts, 2)

		findAccount(t, accounts, account1.ID, func(account beans.AccountWithBalance) {
			assert.False(t, account.ID.Empty())
			assert.Equal(t, account1.Name, account.Name)
			assert.Equal(t, beans.NewAmount(6, 0), account.Balance)
		})
		findAccount(t, accounts, account2.ID, func(account beans.AccountWithBalance) {
			assert.False(t, account.ID.Empty())
			assert.Equal(t, account2.Name, account.Name)
			assert.Equal(t, beans.NewAmount(0, 0), account.Balance)
		})
	})
}
