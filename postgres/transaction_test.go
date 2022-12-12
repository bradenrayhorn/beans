package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	pool, stop := testutils.StartPool(t)
	defer stop()

	transactionRepository := postgres.NewTransactionRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)
	account := makeAccount(t, pool, "account", budgetID)
	categoryGroupID := makeCategoryGroup(t, pool, "group1", budgetID)
	categoryID := makeCategory(t, pool, "category", categoryGroupID, budgetID)

	cleanup := func() {
		pool.Exec(context.Background(), "truncate transactions;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		err := transactionRepository.Create(
			context.Background(),
			&beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(5, 0),
				Date:       beans.NewDate(time.Now()),
				Notes:      beans.NewTransactionNotes("notes"),
			},
		)
		require.Nil(t, err)
	})

	t.Run("can get all", func(t *testing.T) {
		defer cleanup()
		transaction1 := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
			Account:    &account,
		}
		transaction2 := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(7, 0),
			Date:       testutils.NewDate(t, "2022-08-26"),
			Notes:      beans.NewTransactionNotes("my notes"),
			Account:    &account,
		}
		err := transactionRepository.Create(context.Background(), transaction1)
		require.Nil(t, err)
		err = transactionRepository.Create(context.Background(), transaction2)
		require.Nil(t, err)

		transactions, err := transactionRepository.GetForBudget(context.Background(), budgetID)
		assert.Len(t, transactions, 2)
		assert.True(t, reflect.DeepEqual(transactions[0], transaction1))
		assert.True(t, reflect.DeepEqual(transactions[1], transaction2))
	})

	t.Run("can store with empty category", func(t *testing.T) {
		defer cleanup()
		transaction1 := &beans.Transaction{
			ID:        beans.NewBeansID(),
			AccountID: account.ID,
			Amount:    beans.NewAmount(5, 0),
			Date:      testutils.NewDate(t, "2022-08-28"),
		}
		transaction2 := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(7, 0),
			Date:       testutils.NewDate(t, "2022-08-26"),
			CategoryID: testutils.NewEmptyID(),
		}
		require.Nil(t, transactionRepository.Create(context.Background(), transaction1))
		require.Nil(t, transactionRepository.Create(context.Background(), transaction2))

		transactions, err := transactionRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Len(t, transactions, 2)
		assert.True(t, transactions[0].CategoryID.Empty())
		assert.True(t, transactions[1].CategoryID.Empty())
	})
}
