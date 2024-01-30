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
	ds := postgres.NewDataSource(pool)
	factory := testutils.Factory(t, ds)

	transactionRepository := postgres.NewTransactionRepository(pool)

	userID := factory.MakeUser("user")
	budgetID := factory.MakeBudget("budget", userID).ID
	budgetID2 := factory.MakeBudget("budget", userID).ID
	account := factory.MakeAccount("account", budgetID)
	account2 := factory.MakeAccount("account2", budgetID)
	budget2Account1 := factory.MakeAccount("account2", budgetID2)
	categoryGroupID := factory.MakeCategoryGroup("group1", budgetID).ID
	incomeGroup := factory.MakeIncomeCategoryGroup("group2", budgetID)
	budget2IncomeGroup := factory.MakeIncomeCategoryGroup("group2", budgetID2)
	payee := factory.MakePayee("payee", budgetID)
	payee2 := factory.MakePayee("payee2", budgetID)
	categoryID := factory.MakeCategory("category", categoryGroupID, budgetID).ID
	categoryID2 := factory.MakeCategory("category2", categoryGroupID, budgetID).ID
	incomeCategory := factory.MakeCategory("category", incomeGroup.ID, budgetID)
	budget2IncomeCategory := factory.MakeCategory("category", budget2IncomeGroup.ID, budgetID2)

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

	t.Run("can delete", func(t *testing.T) {
		defer cleanup()
		transaction := beans.Transaction{
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
		transaction1 := transaction
		transaction1.ID = beans.NewBeansID()

		transaction2 := transaction
		transaction2.ID = beans.NewBeansID()

		transaction3 := transaction
		transaction2.ID = beans.NewBeansID()

		transaction4 := transaction
		transaction4.ID = beans.NewBeansID()
		transaction4.AccountID = budget2Account1.ID

		require.Nil(t, transactionRepository.Create(context.Background(), &transaction1))
		require.Nil(t, transactionRepository.Create(context.Background(), &transaction2))
		require.Nil(t, transactionRepository.Create(context.Background(), &transaction3))
		require.Nil(t, transactionRepository.Create(context.Background(), &transaction4))

		err := transactionRepository.Delete(context.Background(), budgetID, []beans.ID{transaction1.ID, transaction2.ID, transaction4.ID})
		require.Nil(t, err)

		// transaction1 and transaction2 should be deleted, they are passed in and part of budget 1.
		// transaction3 should not be deleted, it is not passed in.
		// transaction4 should not be deleted, it is passed in but not part of budget 1.
		_, err = transactionRepository.Get(context.Background(), transaction1.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(context.Background(), transaction2.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(context.Background(), transaction3.ID)
		assert.Nil(t, err)

		_, err = transactionRepository.Get(context.Background(), transaction4.ID)
		assert.Nil(t, err)
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
