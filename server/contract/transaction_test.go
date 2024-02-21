package contract_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	t.Parallel()
	pool, ds, factory, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate table users, budgets cascade;")
	}

	transactionRepository := ds.TransactionRepository()
	monthRepository := ds.MonthRepository()
	monthCategoryRepository := ds.MonthCategoryRepository()
	c := contract.NewContracts(ds, inmem.NewSessionRepository()).Transaction

	t.Run("create", func(t *testing.T) {
		t.Run("fields are required", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.Create(context.Background(), auth, beans.TransactionCreateParams{})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		t.Run("cannot create transaction with amount more than 2 decimals", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(10, -3),
					Date:      beans.NewDate(time.Now()),
				},
			}
			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Amount must have at most 2 decimal points.")
		})

		t.Run("can create full", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("group", budget.ID)
			category := factory.MakeCategory("category", group.ID, budget.ID)
			payee := factory.MakePayee("payee", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(1, 2),
					Date:       testutils.NewDate(t, "2022-06-07"),
					Notes:      beans.NewTransactionNotes("My Notes"),
				},
			}

			// transaction was returned
			id, err := c.Create(context.Background(), auth, params)
			require.Nil(t, err)

			// transaction was created
			dbTransactions, err := transactionRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, dbTransactions, 1)

			assert.Equal(t, beans.TransactionWithRelations{
				Transaction: beans.Transaction{
					ID:         id,
					AccountID:  account.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(1, 2),
					Date:       testutils.NewDate(t, "2022-06-07"),
					Notes:      beans.NewTransactionNotes("My Notes"),
				},
				Account:  beans.RelatedAccount{ID: account.ID, Name: account.Name},
				Category: beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}),
				Payee:    beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}),
			}, dbTransactions[0])

			// month was created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 1)
			month := months[0]

			assert.Equal(t, testutils.NewMonthDate(t, "2022-06-01"), month.Date)

			// month category was created
			monthCategories, err := monthCategoryRepository.GetForMonth(context.Background(), month)
			require.Nil(t, err)
			require.Len(t, monthCategories, 1)
			assert.Equal(t, beans.NewAmount(0, 0), monthCategories[0].Amount)
			assert.Equal(t, category.ID, monthCategories[0].CategoryID)
		})

		t.Run("can create minimum", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(1, 2),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			// transaction was returned
			id, err := c.Create(context.Background(), auth, params)
			require.Nil(t, err)

			// transaction was created
			dbTransactions, err := transactionRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, dbTransactions, 1)

			assert.Equal(t, beans.TransactionWithRelations{
				Transaction: beans.Transaction{
					ID:        id,
					AccountID: account.ID,
					Amount:    beans.NewAmount(1, 2),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
				Account: beans.RelatedAccount{ID: account.ID, Name: account.Name},
			}, dbTransactions[0])

			// month was not created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 0)
		})

		t.Run("cannot create with missing account", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: beans.NewID(),
					Amount:    beans.NewAmount(10, 1),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Account ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with account from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			budget2 := factory.MakeBudget("budget", userID)
			account2 := factory.MakeAccount("account", budget2.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account2.ID,
					Amount:    beans.NewAmount(10, 1),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Account ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with non existent category", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					Amount:     beans.NewAmount(10, 1),
					Date:       testutils.NewDate(t, "2022-06-07"),
					CategoryID: beans.NewID(),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with category from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			budget2 := factory.MakeBudget("budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("group", budget2.ID)
			category := factory.MakeCategory("name", group.ID, budget2.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					Amount:     beans.NewAmount(10, 1),
					Date:       testutils.NewDate(t, "2022-06-07"),
					CategoryID: category.ID,
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with non existent payee", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(10, 1),
					Date:      testutils.NewDate(t, "2022-06-07"),
					PayeeID:   beans.NewID(),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with payee from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			budget2 := factory.MakeBudget("budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			payee := factory.MakePayee("payee", budget2.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(10, 1),
					Date:      testutils.NewDate(t, "2022-06-07"),
					PayeeID:   payee.ID,
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("fields are required", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, beans.TransactionUpdateParams{ID: beans.NewID()})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		t.Run("cannot update non existent transaction", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			params := beans.TransactionUpdateParams{
				ID: beans.NewID(),
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Transaction ID")
		})

		t.Run("cannot update transaction for wrong budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)

			budget2 := factory.MakeBudget("budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget2)

			account := factory.MakeAccount("account", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Transaction ID")
		})

		t.Run("can update full", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			account2 := factory.MakeAccount("account2", budget.ID)
			group := factory.MakeCategoryGroup("group", budget.ID)
			category := factory.MakeCategory("category", group.ID, budget.ID)
			payee := factory.MakePayee("payee", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID:  account2.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(6, 0),
					Date:       testutils.NewDate(t, "2022-06-07"),
					Notes:      beans.NewTransactionNotes("My Notes"),
				},
			}

			require.Nil(t, c.Update(context.Background(), auth, params))

			// transaction was updated
			res, err := transactionRepository.Get(context.Background(), budget.ID, transaction.ID)
			require.NoError(t, err)
			assert.Equal(t, beans.Transaction{
				ID:         transaction.ID,
				AccountID:  account2.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(6, 0),
				Date:       testutils.NewDate(t, "2022-06-07"),
				Notes:      beans.NewTransactionNotes("My Notes"),
			}, res)

			// month was created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 1)
			month := months[0]

			assert.Equal(t, testutils.NewMonthDate(t, "2022-06-01"), month.Date)

			// month category was created
			monthCategories, err := monthCategoryRepository.GetForMonth(context.Background(), month)
			require.Nil(t, err)
			require.Len(t, monthCategories, 1)
			assert.Equal(t, beans.NewAmount(0, 0), monthCategories[0].Amount)
			assert.Equal(t, category.ID, monthCategories[0].CategoryID)
		})

		t.Run("can update minimum", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			account2 := factory.MakeAccount("account", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: account2.ID,
					Amount:    beans.NewAmount(6, 0),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			require.Nil(t, c.Update(context.Background(), auth, params))

			// transaction was updated
			res, err := transactionRepository.Get(context.Background(), budget.ID, transaction.ID)
			require.NoError(t, err)
			assert.Equal(t, beans.Transaction{
				ID:        transaction.ID,
				AccountID: account2.ID,
				Amount:    beans.NewAmount(6, 0),
				Date:      testutils.NewDate(t, "2022-06-07"),
			}, res)

			// month was not created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 0)
		})

		t.Run("cannot update with missing account", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := factory.MakeAccount("account", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: beans.NewID(),
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Account ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with account from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := factory.MakeAccount("account", budget.ID)

			budget2 := factory.MakeBudget("budget", userID)
			account2 := factory.MakeAccount("account", budget2.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: account2.ID,
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Account ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with nonexistent category", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					Amount:     beans.NewAmount(5, 0),
					Date:       testutils.NewDate(t, "2023-01-09"),
					CategoryID: beans.NewID(),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with category from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			budget2 := factory.MakeBudget("budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("group", budget2.ID)
			category := factory.MakeCategory("name", group.ID, budget2.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					Amount:     beans.NewAmount(5, 0),
					Date:       testutils.NewDate(t, "2023-01-09"),
					CategoryID: category.ID,
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with nonexistent payee", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
					PayeeID:   beans.NewID(),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with payee from other budget", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			budget2 := factory.MakeBudget("budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			payee := factory.MakePayee("payee", budget2.ID)

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(5, 0),
					Date:      testutils.NewDate(t, "2023-01-09"),
					PayeeID:   payee.ID,
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("can delete", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("group", budget.ID)
			category := factory.MakeCategory("category", group.ID, budget.ID)
			payee := factory.MakePayee("payee", budget.ID)

			transaction := beans.Transaction{
				ID:         beans.NewID(),
				AccountID:  account.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       testutils.NewDate(t, "2023-01-09"),
				Notes:      beans.NewTransactionNotes("hi there"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			require.Nil(t, c.Delete(context.Background(), auth, []beans.ID{transaction.ID}))

			transactions, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, transactions, 0)
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("group", budget.ID)
			category := factory.MakeCategory("category", group.ID, budget.ID)
			payee := factory.MakePayee("payee", budget.ID)

			transaction := beans.Transaction{
				ID:         beans.NewID(),
				AccountID:  account.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       testutils.NewDate(t, "2023-01-09"),
				Notes:      beans.NewTransactionNotes("hi there"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			transactions, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, transactions, 1)

			assert.Equal(t, beans.TransactionWithRelations{
				Transaction: beans.Transaction{
					ID:         transaction.ID,
					AccountID:  account.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(5, 0),
					Date:       testutils.NewDate(t, "2023-01-09"),
					Notes:      beans.NewTransactionNotes("hi there"),
				},
				Account:  beans.RelatedAccount{ID: account.ID, Name: account.Name},
				Category: beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}),
				Payee:    beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}),
			}, transactions[0])
		})
	})
}
