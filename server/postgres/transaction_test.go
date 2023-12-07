package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	transactionRepository := postgres.NewTransactionRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID
	budgetID2 := testutils.MakeBudget(t, pool, "budget", userID).ID
	account := testutils.MakeAccount(t, pool, "account", budgetID)
	account2 := testutils.MakeAccount(t, pool, "account2", budgetID)
	budget2Account1 := testutils.MakeAccount(t, pool, "account2", budgetID2)
	categoryGroupID := testutils.MakeCategoryGroup(t, pool, "group1", budgetID).ID
	incomeGroup := testutils.MakeIncomeCategoryGroup(t, pool, "group2", budgetID)
	budget2IncomeGroup := testutils.MakeIncomeCategoryGroup(t, pool, "group2", budgetID2)
	payee := testutils.MakePayee(t, pool, "payee", budgetID)
	payee2 := testutils.MakePayee(t, pool, "payee2", budgetID)
	categoryID := testutils.MakeCategory(t, pool, "category", categoryGroupID, budgetID).ID
	categoryID2 := testutils.MakeCategory(t, pool, "category2", categoryGroupID, budgetID).ID
	incomeCategory := testutils.MakeCategory(t, pool, "category", incomeGroup.ID, budgetID)
	budget2IncomeCategory := testutils.MakeCategory(t, pool, "category", budget2IncomeGroup.ID, budgetID2)

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate transactions;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		err := transactionRepository.Create(
			context.Background(),
			&beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				CategoryID: categoryID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       beans.NewDate(time.Now()),
				Notes:      beans.NewTransactionNotes("notes"),
			},
		)
		require.Nil(t, err)
	})

	t.Run("cannot get nonexistant", func(t *testing.T) {
		defer cleanup()
		_, err := transactionRepository.Get(context.Background(), beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get", func(t *testing.T) {
		defer cleanup()
		transaction := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			PayeeID:    payee.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
			Account:    account,
		}
		require.Nil(t, transactionRepository.Create(context.Background(), transaction))

		dbTransaction, err := transactionRepository.Get(context.Background(), transaction.ID)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(transaction, dbTransaction))
	})

	t.Run("can update", func(t *testing.T) {
		defer cleanup()
		transaction := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			PayeeID:    payee.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
			Account:    account,
		}
		require.Nil(t, transactionRepository.Create(context.Background(), transaction))

		transaction.AccountID = account2.ID
		transaction.CategoryID = categoryID2
		transaction.PayeeID = payee2.ID
		transaction.Amount = beans.NewAmount(6, 0)
		transaction.Date = testutils.NewDate(t, "2022-08-30")
		transaction.Notes = beans.NewTransactionNotes("notes 5")
		transaction.Account = account2

		require.Nil(t, transactionRepository.Update(context.Background(), transaction))

		dbTransaction, err := transactionRepository.Get(context.Background(), transaction.ID)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(transaction, dbTransaction))
	})

	t.Run("can get all", func(t *testing.T) {
		defer cleanup()
		transaction1 := &beans.Transaction{
			ID:           beans.NewBeansID(),
			AccountID:    account.ID,
			CategoryID:   categoryID,
			PayeeID:      payee.ID,
			Amount:       beans.NewAmount(5, 0),
			Date:         testutils.NewDate(t, "2022-08-28"),
			Notes:        beans.NewTransactionNotes("notes"),
			Account:      account,
			CategoryName: beans.NewNullString("category"),
			PayeeName:    beans.NewNullString("payee"),
		}
		transaction2 := &beans.Transaction{
			ID:           beans.NewBeansID(),
			AccountID:    account.ID,
			CategoryID:   categoryID,
			PayeeID:      payee.ID,
			Amount:       beans.NewAmount(7, 0),
			Date:         testutils.NewDate(t, "2022-08-26"),
			Notes:        beans.NewTransactionNotes("my notes"),
			Account:      account,
			CategoryName: beans.NewNullString("category"),
			PayeeName:    beans.NewNullString("payee"),
		}
		err := transactionRepository.Create(context.Background(), transaction1)
		require.Nil(t, err)
		err = transactionRepository.Create(context.Background(), transaction2)
		require.Nil(t, err)

		transactions, err := transactionRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Len(t, transactions, 2)
		assert.True(t, reflect.DeepEqual(transactions[0], transaction1))
		assert.True(t, reflect.DeepEqual(transactions[1], transaction2))
	})

	t.Run("can store with empty optional fields category", func(t *testing.T) {
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
			PayeeID:    testutils.NewEmptyID(),
		}
		require.Nil(t, transactionRepository.Create(context.Background(), transaction1))
		require.Nil(t, transactionRepository.Create(context.Background(), transaction2))

		transactions, err := transactionRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Len(t, transactions, 2)

		assert.True(t, transactions[0].CategoryID.Empty())
		assert.True(t, transactions[1].CategoryID.Empty())

		assert.True(t, transactions[0].PayeeID.Empty())
		assert.True(t, transactions[1].PayeeID.Empty())
	})

	t.Run("can get income", func(t *testing.T) {
		defer cleanup()
		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(1, 0),
			Date:       testutils.NewDate(t, "2022-09-01"),
			CategoryID: incomeCategory.ID,
		}))
		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(2, 0),
			Date:       testutils.NewDate(t, "2022-08-31"),
			CategoryID: incomeCategory.ID,
		}))
		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(3, 0),
			Date:       testutils.NewDate(t, "2022-08-01"),
			CategoryID: incomeCategory.ID,
		}))
		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(3, 0),
			Date:       testutils.NewDate(t, "2022-07-31"),
			CategoryID: incomeCategory.ID,
		}))

		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-08-31"),
			CategoryID: categoryID,
		}))
		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-07-29"),
			CategoryID: categoryID,
		}))

		require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  budget2Account1.ID,
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-08-15"),
			CategoryID: budget2IncomeCategory.ID,
		}))

		amount, err := transactionRepository.GetIncomeBetween(context.Background(), budgetID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))
		require.Nil(t, err)

		require.Equal(t, beans.NewAmount(5, 0), amount)
	})

	t.Run("can get blank income", func(t *testing.T) {
		defer cleanup()
		amount, err := transactionRepository.GetIncomeBetween(context.Background(), budgetID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))
		require.Nil(t, err)

		require.Equal(t, beans.NewAmount(0, 0), amount)
	})
}
