package contract_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	transactionRepository := postgres.NewTransactionRepository(pool)
	categoryRepository := postgres.NewCategoryRepository(pool)
	accountRepository := postgres.NewAccountRepository(pool)
	monthRepository := postgres.NewMonthRepository(pool)
	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)
	payeeRepository := postgres.NewPayeeRepository(pool)
	c := contract.NewTransactionContract(
		transactionRepository,
		accountRepository,
		categoryRepository,
		monthCategoryRepository,
		monthRepository,
		payeeRepository,
	)

	t.Run("create", func(t *testing.T) {
		t.Run("fields are required", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.Create(context.Background(), auth, beans.TransactionCreateParams{})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		t.Run("cannot create transaction with amount more than 2 decimals", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "group", budget.ID)
			category := testutils.MakeCategory(t, pool, "category", group.ID, budget.ID)
			payee := testutils.MakePayee(t, pool, "payee", budget.ID)

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
			transaction, err := c.Create(context.Background(), auth, params)
			require.Nil(t, err)
			require.Equal(t, params.AccountID, transaction.AccountID)
			require.Equal(t, params.CategoryID, transaction.CategoryID)
			require.Equal(t, params.PayeeID, transaction.PayeeID)
			require.Equal(t, params.Amount, transaction.Amount)
			require.Equal(t, params.Date, transaction.Date)
			require.Equal(t, params.Notes, transaction.Notes)
			assert.True(t, reflect.DeepEqual(account, transaction.Account))

			// transaction was created
			dbTransactions, err := transactionRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, dbTransactions, 1)

			transaction.CategoryName = beans.NewNullString("category")
			transaction.PayeeName = beans.NewNullString("payee")
			assert.True(t, reflect.DeepEqual(transaction, dbTransactions[0]))

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(1, 2),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			// transaction was returned
			transaction, err := c.Create(context.Background(), auth, params)
			require.Nil(t, err)
			require.Equal(t, params.AccountID, transaction.AccountID)
			require.Equal(t, params.CategoryID, transaction.CategoryID)
			require.Equal(t, params.PayeeID, transaction.PayeeID)
			require.Equal(t, params.Amount, transaction.Amount)
			require.Equal(t, params.Date, transaction.Date)
			require.Equal(t, params.Notes, transaction.Notes)
			assert.True(t, reflect.DeepEqual(account, transaction.Account))

			// transaction was created
			dbTransactions, err := transactionRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, dbTransactions, 1)

			assert.True(t, reflect.DeepEqual(transaction, dbTransactions[0]))

			// month was not created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 0)
		})

		t.Run("cannot create with missing account", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			budget2 := testutils.MakeBudget(t, pool, "budget", userID)
			account2 := testutils.MakeAccount(t, pool, "account", budget2.ID)

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					Amount:     beans.NewAmount(10, 1),
					Date:       testutils.NewDate(t, "2022-06-07"),
					CategoryID: beans.NewBeansID(),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with category from other budget", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "group", budget2.ID)
			category := testutils.MakeCategory(t, pool, "name", group.ID, budget2.ID)

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(10, 1),
					Date:      testutils.NewDate(t, "2022-06-07"),
					PayeeID:   beans.NewBeansID(),
				},
			}

			_, err := c.Create(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with payee from other budget", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			payee := testutils.MakePayee(t, pool, "payee", budget2.ID)

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, beans.TransactionUpdateParams{ID: beans.NewBeansID()})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		t.Run("cannot update non existent transaction", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			params := beans.TransactionUpdateParams{
				ID: beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)

			budget2 := testutils.MakeBudget(t, pool, "budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget2)

			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			account2 := testutils.MakeAccount(t, pool, "account2", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "group", budget.ID)
			category := testutils.MakeCategory(t, pool, "category", group.ID, budget.ID)
			payee := testutils.MakePayee(t, pool, "payee", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),

				Account: account,
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

			transaction.AccountID = account2.ID
			transaction.CategoryID = category.ID
			transaction.PayeeID = payee.ID
			transaction.Amount = beans.NewAmount(6, 0)
			transaction.Date = testutils.NewDate(t, "2022-06-07")
			transaction.Notes = beans.NewTransactionNotes("My Notes")

			transaction.Account = account2

			require.Nil(t, c.Update(context.Background(), auth, params))

			// transaction was updated
			dbTransaction, err := transactionRepository.Get(context.Background(), transaction.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(transaction, dbTransaction))

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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			account2 := testutils.MakeAccount(t, pool, "account", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),

				Account: account,
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

			transaction.AccountID = account2.ID
			transaction.Amount = beans.NewAmount(6, 0)
			transaction.Date = testutils.NewDate(t, "2022-06-07")

			transaction.Account = account2

			require.Nil(t, c.Update(context.Background(), auth, params))

			// transaction was updated
			dbTransaction, err := transactionRepository.Get(context.Background(), transaction.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(transaction, dbTransaction))

			// month was not created
			months, err := monthRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, months, 0)
		})

		t.Run("cannot update with missing account", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2023-01-09"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			params := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					AccountID: beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			budget2 := testutils.MakeBudget(t, pool, "budget", userID)
			account2 := testutils.MakeAccount(t, pool, "account", budget2.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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
					CategoryID: beans.NewBeansID(),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Category ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with category from other budget", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "group", budget2.ID)
			category := testutils.MakeCategory(t, pool, "name", group.ID, budget2.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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
					PayeeID:   beans.NewBeansID(),
				},
			}

			err := c.Update(context.Background(), auth, params)
			testutils.AssertError(t, err, "Invalid Payee ID")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot update with payee from other budget", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "budget2", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			payee := testutils.MakePayee(t, pool, "payee", budget2.ID)

			transaction := &beans.Transaction{
				ID:        beans.NewBeansID(),
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

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "group", budget.ID)
			category := testutils.MakeCategory(t, pool, "category", group.ID, budget.ID)
			payee := testutils.MakePayee(t, pool, "payee", budget.ID)

			transaction := &beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       testutils.NewDate(t, "2023-01-09"),
				Notes:      beans.NewTransactionNotes("hi there"),

				Account:      account,
				CategoryName: beans.NewNullString("category"),
				PayeeName:    beans.NewNullString("payee"),
			}
			require.Nil(t, transactionRepository.Create(context.Background(), transaction))

			transactions, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, transactions, 1)

			assert.True(t, reflect.DeepEqual(transaction, transactions[0]))
		})
	})
}