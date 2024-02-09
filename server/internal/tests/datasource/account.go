package datasource

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository(t *testing.T, ds beans.DataSource) {
	factory := testutils.Factory(t, ds)
	accountRepository := ds.AccountRepository()
	ctx := context.Background()

	t.Run("can create and get account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		accountID := beans.NewBeansID()
		err := accountRepository.Create(
			ctx,
			accountID,
			"Account1",
			budget.ID,
		)
		require.Nil(t, err)

		account, err := accountRepository.Get(context.Background(), budget.ID, accountID)
		require.Nil(t, err)
		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, "Account1", string(account.Name))
		assert.Equal(t, budget.ID, account.BudgetID)
	})

	t.Run("cannot create duplicate account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		accountID := beans.NewBeansID()

		err := accountRepository.Create(ctx, accountID, "Account1", budget.ID)
		require.Nil(t, err)

		err = accountRepository.Create(ctx, accountID, "Account1", budget.ID)
		require.NotNil(t, err)
	})

	t.Run("cannot get fictitious account", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		accountID := beans.NewBeansID()

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

	t.Run("can get accounts for budget", func(t *testing.T) {
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

		res, err := accountRepository.GetForBudget(ctx, budget1.ID)
		require.Nil(t, err)

		// accounts 1 and 2 should be in the response with a balance
		expectedAccounts := []beans.AccountWithBalance{
			{Account: account1, Balance: beans.NewAmount(0, 0)},
			{Account: account2, Balance: beans.NewAmount(2, 0)},
		}

		testutils.IsEqualInAnyOrder(t, res, expectedAccounts, testutils.CmpAccountWithBalance)
	})
}
