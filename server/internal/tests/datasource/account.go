package datasource

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAccount(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)
	accountRepository := ds.AccountRepository()
	ctx := context.Background()

	t.Run("can create and get account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		account := beans.Account{
			ID:        beans.NewID(),
			BudgetID:  budget.ID,
			Name:      beans.Name("Account1"),
			OffBudget: true,
		}
		err := accountRepository.Create(ctx, account)
		require.NoError(t, err)

		res, err := accountRepository.Get(context.Background(), budget.ID, account.ID)
		require.NoError(t, err)
		assert.Equal(t, account, res)
	})

	t.Run("cannot create duplicate account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		account := beans.Account{
			ID:       beans.NewID(),
			BudgetID: budget.ID,
			Name:     beans.Name("Account1"),
		}

		err := accountRepository.Create(ctx, account)
		require.NoError(t, err)

		err = accountRepository.Create(ctx, account)
		require.NotNil(t, err)
	})

	t.Run("cannot get fictitious account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		accountID := beans.NewID()

		_, err := accountRepository.Get(ctx, accountID, budget.ID)

		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("cannot get account for other budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})

		_, err := accountRepository.Get(ctx, account.ID, budget2.ID)

		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("get with balance", func(t *testing.T) {

		t.Run("can get accounts with balance", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			// accounts 1 and 2 should be in the response
			account1 := factory.Account(beans.Account{BudgetID: budget1.ID})
			account2 := factory.Account(beans.Account{BudgetID: budget1.ID})
			account3 := factory.Account(beans.Account{BudgetID: budget2.ID})

			// account 2 should have a balance of $2
			factory.Transaction(budget1.ID, beans.Transaction{
				AccountID: account2.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2022-05-20"),
			})
			factory.Transaction(budget1.ID, beans.Transaction{
				AccountID: account2.ID,
				Amount:    beans.NewAmount(-3, 0),
				Date:      testutils.NewDate(t, "1945-05-20"),
			})
			// account 3's transactions should not impact the balances
			factory.Transaction(budget2.ID, beans.Transaction{
				AccountID: account3.ID,
				Amount:    beans.NewAmount(-3, 0),
				Date:      testutils.NewDate(t, "1945-05-23"),
			})

			res, err := accountRepository.GetWithBalance(ctx, budget1.ID)
			require.NoError(t, err)

			// accounts 1 and 2 should be in the response with a balance
			expectedAccounts := []beans.AccountWithBalance{
				{Account: account1, Balance: beans.NewAmount(0, 0)},
				{Account: account2, Balance: beans.NewAmount(2, 0)},
			}

			assert.ElementsMatch(t, expectedAccounts, res)
		})

		t.Run("includes off budget accounts", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})

			res, err := accountRepository.GetWithBalance(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, account.ID, res[0].ID)
		})

		t.Run("excludes split parent from sum", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})
			parent := factory.Transaction(budget.ID, beans.Transaction{
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2022-05-20"),
				IsSplit:   true,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2022-05-20"),
				SplitID:   parent.ID,
			})

			res, err := accountRepository.GetWithBalance(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, beans.NewAmount(5, 0), res[0].Balance)
		})
	})

	t.Run("get transactable", func(t *testing.T) {

		t.Run("filters by budget", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			factory.Account(beans.Account{BudgetID: budget1.ID})

			// budget 2 should have no accounts
			res, err := accountRepository.GetTransactable(ctx, budget2.ID)
			require.NoError(t, err)

			assert.Equal(t, 0, len(res))
		})

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})

			// the account should be returned
			res, err := accountRepository.GetTransactable(ctx, budget.ID)
			require.NoError(t, err)

			assert.Equal(t, 1, len(res))
			assert.Equal(t, account, res[0])
		})
	})

}
