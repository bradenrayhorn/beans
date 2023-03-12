package contract_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/contract"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonth(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	categoryRepository := postgres.NewCategoryRepository(pool)
	monthRepository := postgres.NewMonthRepository(pool)
	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)
	transactionRepository := postgres.NewTransactionRepository(pool)
	txManager := postgres.NewTxManager(pool)
	c := contract.NewMonthContract(categoryRepository, monthRepository, monthCategoryRepository, transactionRepository, txManager)

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get non existant month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, _, _, err := c.Get(context.Background(), auth, beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("must have access to month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "Budget2", userID)
			month := testutils.MakeMonth(t, pool, budget2.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, _, _, err := c.Get(context.Background(), auth, month.ID)
			testutils.AssertErrorCode(t, err, beans.EFORBIDDEN)
		})

		t.Run("can get month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			monthApril := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-04-01"))
			monthApril.Carryover = beans.NewAmount(67, -1)
			require.Nil(t, monthRepository.Update(context.Background(), monthApril))

			monthMay := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))
			monthMay.Carryover = beans.NewAmount(4, -1)
			require.Nil(t, monthRepository.Update(context.Background(), monthMay))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := testutils.MakeAccount(t, pool, "account", budget.ID)
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			incomeGroup := testutils.MakeIncomeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)
			incomeCategory := testutils.MakeCategory(t, pool, "Income", incomeGroup.ID, budget.ID)
			monthCategory := testutils.MakeMonthCategory(t, pool, monthMay.ID, category.ID, beans.NewAmount(34, -1))
			testutils.MakeMonthCategory(t, pool, monthApril.ID, category.ID, beans.NewAmount(34, -1))

			require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(6, 0),
				Date:       testutils.NewDate(t, "2022-03-01"),
				CategoryID: incomeCategory.ID,
			}))
			require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(9, 0),
				Date:       testutils.NewDate(t, "2022-05-01"),
				CategoryID: incomeCategory.ID,
			}))
			require.Nil(t, transactionRepository.Create(context.Background(), &beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-06-01"),
				CategoryID: incomeCategory.ID,
			}))

			dbMonth, dbCategories, available, err := c.Get(context.Background(), auth, monthMay.ID)
			require.Nil(t, err)

			monthMay.Income = beans.NewAmount(9, 0)
			monthMay.Assigned = beans.NewAmount(34, -1)
			monthMay.CarriedOver = beans.NewAmount(67, -1)
			assert.True(t, reflect.DeepEqual(monthMay, dbMonth))
			require.Len(t, dbCategories, 1)

			monthCategory.Activity = beans.NewAmount(0, 0)
			monthCategory.Available = beans.NewAmount(68, -1)
			assert.True(t, reflect.DeepEqual(monthCategory, dbCategories[0]))

			assert.Equal(t, beans.NewAmount(119, -1), available)
		})
	})

	t.Run("create", func(t *testing.T) {
		t.Run("creates new month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			date := testutils.NewMonthDate(t, "2022-05-01")

			month, err := c.CreateMonth(context.Background(), auth, date)
			require.Nil(t, err)

			// month was returned
			assert.False(t, month.ID.Empty())
			assert.Equal(t, budget.ID, month.BudgetID)
			assert.Equal(t, date, month.Date)

			// month was saved
			dbMonth, err := monthRepository.Get(context.Background(), month.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(month, dbMonth))
		})

		t.Run("uses existing month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			returnedMonth, err := c.CreateMonth(context.Background(), auth, month.Date)
			require.Nil(t, err)

			// month was returned
			assert.True(t, reflect.DeepEqual(month, returnedMonth))
		})

		t.Run("creates existing month categories when creating month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Electric", group.ID, budget.ID)

			date := testutils.NewMonthDate(t, "2022-05-01")

			month, err := c.CreateMonth(context.Background(), auth, date)
			require.Nil(t, err)

			monthCategories, err := monthCategoryRepository.GetForMonth(context.Background(), month)
			require.Nil(t, err)
			require.Len(t, monthCategories, 1)
			require.Equal(t, category.ID, monthCategories[0].CategoryID)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("cannot update non existant month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, beans.NewBeansID(), beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("must have access to month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "Budget2", userID)
			month := testutils.MakeMonth(t, pool, budget2.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.EFORBIDDEN)
		})

		t.Run("cannot add negative carryover", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(-5, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot add blank carryover", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewEmptyAmount())
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can update carryover", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			res, err := monthRepository.Get(context.Background(), month.ID)
			require.Nil(t, err)
			assert.Equal(t, beans.NewAmount(5, 0), res.Carryover)
		})

		t.Run("can update carryover to zero", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(0, 0))
			require.Nil(t, err)

			res, err := monthRepository.Get(context.Background(), month.ID)
			require.Nil(t, err)
			assert.Equal(t, beans.NewAmount(0, 0), res.Carryover)
		})
	})

	t.Run("set category amount", func(t *testing.T) {
		t.Run("amount must be not be zero", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID}

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, beans.NewBeansID(), beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("amount must be not be negative", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID}

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, beans.NewBeansID(), beans.NewAmount(-5, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("must have access to month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			budget2 := testutils.MakeBudget(t, pool, "Budget2", userID)
			month := testutils.MakeMonth(t, pool, budget2.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			testutils.AssertErrorCode(t, err, beans.EFORBIDDEN)
		})

		t.Run("creates new month category", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			monthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, month.ID, category.ID)
			require.Nil(t, err)

			assert.True(t, reflect.DeepEqual(
				monthCategory,
				&beans.MonthCategory{
					ID:         monthCategory.ID,
					CategoryID: category.ID,
					MonthID:    month.ID,
					Amount:     beans.NewAmount(5, 0),
				},
			))
		})

		t.Run("uses existing month category", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)
			monthCategory := testutils.MakeMonthCategory(t, pool, month.ID, category.ID, beans.NewAmount(4, 0))

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			dbMonthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, month.ID, category.ID)
			require.Nil(t, err)

			assert.True(t, reflect.DeepEqual(
				dbMonthCategory,
				&beans.MonthCategory{
					ID:         monthCategory.ID,
					CategoryID: category.ID,
					MonthID:    month.ID,
					Amount:     beans.NewAmount(5, 0),
				},
			))
		})

	})
}
