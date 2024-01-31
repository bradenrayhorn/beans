package datasource

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository(t *testing.T, ds beans.DataSource) {
	factory := testutils.Factory(t, ds)

	transactionRepository := ds.TransactionRepository()
	ctx := context.Background()

	t.Run("can create", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		err := transactionRepository.Create(
			ctx,
			&beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       beans.NewDate(time.Now()),
				Notes:      beans.NewTransactionNotes("notes"),
			},
		)
		require.Nil(t, err)
	})

	t.Run("can create with empty optional fields", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})

		transaction1 := &beans.Transaction{
			ID:        beans.NewBeansID(),
			AccountID: account.ID,
			Amount:    beans.NewAmount(5, 0),
			Date:      testutils.NewDate(t, "2022-08-28"),
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction1))

		transactions, err := transactionRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		assert.Len(t, transactions, 1)

		assert.True(t, transactions[0].CategoryID.Empty())
		assert.True(t, transactions[0].PayeeID.Empty())
		assert.True(t, transactions[0].Notes.Empty())
	})

	t.Run("cannot get nonexistant", func(t *testing.T) {
		_, err := transactionRepository.Get(ctx, beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: category.ID,
			PayeeID:    payee.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction))

		res, err := transactionRepository.Get(ctx, transaction.ID)
		require.Nil(t, err)

		// Account should have been attached
		transaction.Account = &account
		assert.True(t, reflect.DeepEqual(transaction, res))
	})

	t.Run("can update", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account1 := factory.Account(beans.Account{BudgetID: budget.ID})
		account2 := factory.Account(beans.Account{BudgetID: budget.ID})
		payee1 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		payee2 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category1 := factory.Category(beans.Category{BudgetID: budget.ID})
		category2 := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account1.ID,
			CategoryID: category1.ID,
			PayeeID:    payee1.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
			Account:    &account1,
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction))

		transaction.AccountID = account2.ID
		transaction.CategoryID = category2.ID
		transaction.PayeeID = payee2.ID
		transaction.Amount = beans.NewAmount(6, 0)
		transaction.Date = testutils.NewDate(t, "2022-08-30")
		transaction.Notes = beans.NewTransactionNotes("notes 5")
		transaction.Account = &account2

		require.Nil(t, transactionRepository.Update(ctx, transaction))

		res, err := transactionRepository.Get(ctx, transaction.ID)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(transaction, res))
	})

	t.Run("can delete", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		transaction1 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction2 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction3 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction4 := factory.Transaction(budget2.ID, beans.Transaction{})

		err := transactionRepository.Delete(ctx, budget1.ID, []beans.ID{transaction1.ID, transaction2.ID, transaction4.ID})
		require.Nil(t, err)

		// transaction1 and transaction2 should be deleted, they are passed in and part of budget 1.
		// transaction3 should not be deleted, it is not passed in.
		// transaction4 should not be deleted, it is passed in but not part of budget 1.
		_, err = transactionRepository.Get(ctx, transaction1.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(ctx, transaction2.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(ctx, transaction3.ID)
		assert.Nil(t, err)

		_, err = transactionRepository.Get(ctx, transaction4.ID)
		assert.Nil(t, err)
	})

	t.Run("can get all", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		account := factory.Account(beans.Account{BudgetID: budget1.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget1.ID})
		category := factory.Category(beans.Category{BudgetID: budget1.ID})

		transaction1 := factory.Transaction(budget1.ID, beans.Transaction{
			AccountID:  account.ID,
			PayeeID:    payee.ID,
			CategoryID: category.ID,
		})
		// this transaction should not be included
		factory.Transaction(budget2.ID, beans.Transaction{})

		transactions, err := transactionRepository.GetForBudget(ctx, budget1.ID)
		require.Nil(t, err)
		assert.Len(t, transactions, 1)

		// Account, CategoryName, PayeeName, should be loaded
		transaction1.Account = &account
		transaction1.CategoryName = beans.NewNullString(string(category.Name))
		transaction1.PayeeName = beans.NewNullString(string(payee.Name))
		fmt.Printf("%v", transactions[0])
		fmt.Printf("%v", &transaction1)
		fmt.Printf("%v", transactions[0].Amount)
		fmt.Printf("%v", transaction1.Amount)
		fmt.Printf("%v", transactions[0].Amount.Exponent())
		fmt.Printf("%v", transaction1.Amount.Exponent())
		assert.True(t, reflect.DeepEqual(transactions[0], &transaction1))
	})

	t.Run("can get income", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		incomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget.ID, IsIncome: true})
		incomeCategory := factory.Category(beans.Category{BudgetID: budget.ID, GroupID: incomeGroup.ID})
		otherCategory := factory.Category(beans.Category{BudgetID: budget.ID})

		budget2, _ := factory.MakeBudgetAndUser()
		budget2IncomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget2.ID, IsIncome: true})
		budget2IncomeCategory := factory.Category(beans.Category{BudgetID: budget2.ID, GroupID: budget2IncomeGroup.ID})

		// Earned $1 in September
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(1, 0),
			Date:       testutils.NewDate(t, "2022-09-01"),
			CategoryID: incomeCategory.ID,
		})
		// Earned $2 in August
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(2, 0),
			Date:       testutils.NewDate(t, "2022-08-31"),
			CategoryID: incomeCategory.ID,
		})
		// Earned $3 in August
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(3, 0),
			Date:       testutils.NewDate(t, "2022-08-01"),
			CategoryID: incomeCategory.ID,
		})
		// Earned $3 in July
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(3, 0),
			Date:       testutils.NewDate(t, "2022-07-31"),
			CategoryID: incomeCategory.ID,
		})

		// Spent $99 in August
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-08-31"),
			CategoryID: otherCategory.ID,
		})
		// Spent $99 in July
		factory.Transaction(budget.ID, beans.Transaction{
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-07-29"),
			CategoryID: otherCategory.ID,
		})

		// Budget 2, earned $99 in August
		factory.Transaction(budget2.ID, beans.Transaction{
			Amount:     beans.NewAmount(99, 0),
			Date:       testutils.NewDate(t, "2022-08-15"),
			CategoryID: budget2IncomeCategory.ID,
		})

		amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))
		require.Nil(t, err)

		// August earnings for budget 1 = $5
		require.Equal(t, beans.NewAmount(5, 0), amount)
	})

	t.Run("can get blank income", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))

		require.Nil(t, err)

		require.Equal(t, beans.NewAmount(0, 0), amount)
	})
}
