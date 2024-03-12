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

			_, err := interactor.AccountCreate(t, c.ctx, beans.AccountCreate{
				Name: beans.Name(""),
			})
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create account
			accountID, err := interactor.AccountCreate(t, c.ctx, beans.AccountCreate{
				Name: beans.Name("New Account"),
			})
			require.NoError(t, err)

			// check if account was saved properly
			account, err := interactor.AccountGet(t, c.ctx, accountID)
			require.NoError(t, err)

			assert.False(t, account.ID.Empty())
			assert.Equal(t, beans.Name("New Account"), account.Name)
			assert.Equal(t, false, account.OffBudget)
		})

		t.Run("can create and get an off budget account", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create account
			accountID, err := interactor.AccountCreate(t, c.ctx, beans.AccountCreate{
				Name:      beans.Name("New Account"),
				OffBudget: true,
			})
			require.NoError(t, err)

			// check if account was saved properly
			account, err := interactor.AccountGet(t, c.ctx, accountID)
			require.NoError(t, err)

			assert.False(t, account.ID.Empty())
			assert.Equal(t, beans.Name("New Account"), account.Name)
			assert.Equal(t, true, account.OffBudget)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get a non-existent account", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.AccountGet(t, c.ctx, beans.NewID())
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

		t.Run("can", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account1 := c.Account(AccountOpts{})
			account2 := c.Account(AccountOpts{OffBudget: true})

			// add a +$6 transaction to account 1
			categoryGroup := c.CategoryGroup(CategoryGroupOpts{})
			category := c.Category(CategoryOpts{Group: categoryGroup})
			c.Transaction(TransactionOpts{
				Account:  account1,
				Category: category,
				Amount:   "6",
			})

			// list accounts, check if accounts are proper
			accounts, err := interactor.AccountList(t, c.ctx)
			require.NoError(t, err)
			require.Len(t, accounts, 2)

			findAccountWithBalance(t, accounts, account1.ID, func(account beans.AccountWithBalance) {
				assert.False(t, account.ID.Empty())
				assert.Equal(t, account1.Name, account.Name)
				assert.Equal(t, beans.NewAmount(6, 0), account.Balance)
				assert.Equal(t, false, account.OffBudget)
			})
			findAccountWithBalance(t, accounts, account2.ID, func(account beans.AccountWithBalance) {
				assert.False(t, account.ID.Empty())
				assert.Equal(t, account2.Name, account.Name)
				assert.Equal(t, beans.NewAmount(0, 0), account.Balance)
				assert.Equal(t, true, account.OffBudget)
			})
		})
	})

	t.Run("list transactable", func(t *testing.T) {

		t.Run("can", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{OffBudget: true})

			accounts, err := interactor.AccountListTransactable(t, c.ctx)
			require.NoError(t, err)
			require.Equal(t, 1, len(accounts))

			findAccount(t, accounts, account.ID, func(it beans.Account) {
				assert.Equal(t, false, it.ID.Empty())
				assert.Equal(t, account.Name, it.Name)
				assert.Equal(t, true, it.OffBudget)
			})
		})
	})
}
