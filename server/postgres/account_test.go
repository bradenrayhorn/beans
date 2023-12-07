package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	accountRepository := postgres.NewAccountRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID
	budgetID2 := testutils.MakeBudget(t, pool, "budget", userID).ID

	categoryGroup := testutils.MakeCategoryGroup(t, pool, "group", budgetID)
	category1 := testutils.MakeCategory(t, pool, "cat1", categoryGroup.ID, budgetID)
	category2 := testutils.MakeCategory(t, pool, "cat2", categoryGroup.ID, budgetID)

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate accounts cascade;")
	}

	t.Run("can create and get account", func(t *testing.T) {
		defer cleanup()
		accountID := beans.NewBeansID()
		err := accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		require.Nil(t, err)

		account, err := accountRepository.Get(context.Background(), accountID)
		require.Nil(t, err)
		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, "Account1", string(account.Name))
		assert.Equal(t, budgetID, account.BudgetID)
	})

	t.Run("cannot create duplicate account", func(t *testing.T) {
		defer cleanup()
		accountID := beans.NewBeansID()
		err := accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		require.Nil(t, err)

		err = accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		assertPgError(t, pgerrcode.UniqueViolation, err)
	})

	t.Run("cannot get fictitious account", func(t *testing.T) {
		defer cleanup()
		accountID := beans.NewBeansID()
		_, err := accountRepository.Get(context.Background(), accountID)
		require.NotNil(t, err)
		var beansError beans.Error
		require.ErrorAs(t, err, &beansError)
		code, _ := beansError.BeansError()
		assert.Equal(t, beans.ENOTFOUND, code)
	})

	t.Run("can get accounts for budget", func(t *testing.T) {
		defer cleanup()

		account1 := beans.Account{ID: beans.NewBeansID(), BudgetID: budgetID, Name: "Account", Balance: beans.NewAmount(0, 0)}
		account2 := beans.Account{ID: beans.NewBeansID(), BudgetID: budgetID, Name: "Account", Balance: beans.NewAmount(2, 0)}
		account3 := beans.Account{ID: beans.NewBeansID(), BudgetID: budgetID2, Name: "Account", Balance: beans.NewAmount(0, 0)}

		require.Nil(t, accountRepository.Create(context.Background(), account1.ID, account1.Name, account1.BudgetID))
		require.Nil(t, accountRepository.Create(context.Background(), account2.ID, account2.Name, account2.BudgetID))
		require.Nil(t, accountRepository.Create(context.Background(), account3.ID, account3.Name, account3.BudgetID))

		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account2.ID,
			CategoryID: category1.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-05-20"),
		})
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account2.ID,
			CategoryID: category2.ID,
			Amount:     beans.NewAmount(-3, 0),
			Date:       testutils.NewDate(t, "1945-05-20"),
		})

		res, err := accountRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		require.Len(t, res, 2)

		res1 := findResult(res, func(a *beans.Account) bool { return a.ID == account1.ID })
		res2 := findResult(res, func(a *beans.Account) bool { return a.ID == account2.ID })

		assert.True(t, reflect.DeepEqual(res1, &account1))
		assert.True(t, reflect.DeepEqual(res2, &account2))
	})
}
