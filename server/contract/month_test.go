package contract_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonth(t *testing.T) {
	t.Parallel()
	pool, ds, factory, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate table users, budgets cascade;")
	}

	monthRepository := ds.MonthRepository()
	monthCategoryRepository := ds.MonthCategoryRepository()
	transactionRepository := ds.TransactionRepository()
	c := contract.NewContracts(ds, inmem.NewSessionRepository()).Month

	t.Run("get or create", func(t *testing.T) {
		t.Run("creates new month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			date := testutils.NewMonthDate(t, "2022-05-01")

			month, err := c.GetOrCreate(context.Background(), auth, date)
			require.Nil(t, err)

			// month was returned
			assert.Equal(t, budget.ID, month.BudgetID)
			assert.Equal(t, date, month.Date)

			// month was saved
			res, err := monthRepository.Get(context.Background(), budget.ID, month.ID)
			require.NoError(t, err)

			assert.Equal(t, month.Month, res)
		})

		t.Run("uses existing month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			res, err := c.GetOrCreate(context.Background(), auth, month.Date)
			require.NoError(t, err)

			// month was returned
			assert.Equal(t, beans.MonthWithDetails{
				Month:       month,
				CarriedOver: beans.NewAmount(0, 0),
				Income:      beans.NewAmount(0, 0),
				Assigned:    beans.NewAmount(0, 0),
				Budgetable:  beans.NewAmount(0, 0),
				Categories:  []beans.MonthCategoryWithDetails{},
			}, res)
		})

		t.Run("does not use other budget's existing month", func(t *testing.T) {
			defer cleanup()

			date := testutils.NewDate(t, "2022-05-01")

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, date)

			budget2 := factory.MakeBudget("Budget2", userID)
			factory.MakeMonth(budget2.ID, date)

			auth := testutils.BudgetAuthContext(t, userID, budget)

			res, err := c.GetOrCreate(context.Background(), auth, beans.NewMonthDate(date))
			require.Nil(t, err)
			require.Equal(t, month.ID, res.ID)
		})

		t.Run("creates existing month categories when creating month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := factory.MakeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Electric", group.ID, budget.ID)

			date := testutils.NewMonthDate(t, "2022-05-01")

			month, err := c.GetOrCreate(context.Background(), auth, date)
			require.Nil(t, err)

			monthCategories, err := monthCategoryRepository.GetForMonth(context.Background(), month.Month)
			require.Nil(t, err)
			require.Len(t, monthCategories, 1)
			require.Equal(t, category.ID, monthCategories[0].CategoryID)
		})

		t.Run("can get month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			monthApril := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-04-01"))
			monthApril.Carryover = beans.NewAmount(67, -1)
			require.Nil(t, monthRepository.Update(context.Background(), monthApril))

			monthMay := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))
			monthMay.Carryover = beans.NewAmount(4, -1)
			require.Nil(t, monthRepository.Update(context.Background(), monthMay))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			account := factory.MakeAccount("account", budget.ID)
			group := factory.MakeCategoryGroup("Group", budget.ID)
			incomeGroup := factory.MakeIncomeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Category", group.ID, budget.ID)
			incomeCategory := factory.MakeCategory("Income", incomeGroup.ID, budget.ID)
			monthCategory := factory.MakeMonthCategory(monthMay.ID, category.ID, beans.NewAmount(34, -1))
			factory.MakeMonthCategory(monthApril.ID, category.ID, beans.NewAmount(34, -1))

			require.Nil(t, transactionRepository.Create(context.Background(), beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(6, 0),
				Date:       testutils.NewDate(t, "2022-03-01"),
				CategoryID: incomeCategory.ID,
			}))
			require.Nil(t, transactionRepository.Create(context.Background(), beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(9, 0),
				Date:       testutils.NewDate(t, "2022-05-01"),
				CategoryID: incomeCategory.ID,
			}))
			require.Nil(t, transactionRepository.Create(context.Background(), beans.Transaction{
				ID:         beans.NewBeansID(),
				AccountID:  account.ID,
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-06-01"),
				CategoryID: incomeCategory.ID,
			}))

			res, err := c.GetOrCreate(context.Background(), auth, monthMay.Date)
			require.Nil(t, err)

			assert.Equal(t, beans.NewAmount(9, 0), res.Income)
			assert.Equal(t, beans.NewAmount(34, -1), res.Assigned)
			assert.Equal(t, beans.NewAmount(67, -1), res.CarriedOver)
			assert.Equal(t, beans.NewAmount(119, -1), res.Budgetable)
			assert.Equal(t, monthMay, res.Month)

			require.Len(t, res.Categories, 2)

			var dbExpenseCategory beans.MonthCategoryWithDetails
			var dbIncomeCategory beans.MonthCategoryWithDetails
			for _, c := range res.Categories {
				if c.CategoryID == category.ID {
					dbExpenseCategory = c
				}
				if c.CategoryID == incomeCategory.ID {
					dbIncomeCategory = c
				}
			}

			assert.Equal(t, beans.MonthCategoryWithDetails{
				MonthCategory: monthCategory,
				Activity:      beans.NewAmount(0, 0),
				Available:     beans.NewAmount(68, -1),
			}, dbExpenseCategory)

			assert.Equal(t, beans.NewAmount(9, 0), dbIncomeCategory.Activity)
			assert.Equal(t, beans.NewAmount(15, 0), dbIncomeCategory.Available)
			assert.Equal(t, beans.NewAmount(0, 0), dbIncomeCategory.Amount)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("cannot update non existant month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, beans.NewBeansID(), beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("must have access to month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			budget2 := factory.MakeBudget("Budget2", userID)
			month := factory.MakeMonth(budget2.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot add negative carryover", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(-5, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot add blank carryover", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))
			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewEmptyAmount())
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can update carryover", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			res, err := monthRepository.Get(context.Background(), budget.ID, month.ID)
			require.Nil(t, err)
			assert.Equal(t, beans.NewAmount(5, 0), res.Carryover)
		})

		t.Run("can update carryover to zero", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			err := c.Update(context.Background(), auth, month.ID, beans.NewAmount(0, 0))
			require.Nil(t, err)

			res, err := monthRepository.Get(context.Background(), budget.ID, month.ID)
			require.Nil(t, err)
			assert.Equal(t, beans.NewAmount(0, 0), res.Carryover)
		})
	})

	t.Run("set category amount", func(t *testing.T) {
		t.Run("amount must be not be zero", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID}

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, beans.NewBeansID(), beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("amount must be not be negative", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID}

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, beans.NewBeansID(), beans.NewAmount(-5, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("must have access to month", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			budget2 := factory.MakeBudget("Budget2", userID)
			month := factory.MakeMonth(budget2.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := factory.MakeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Category", group.ID, budget.ID)

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("creates new month category", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := factory.MakeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Category", group.ID, budget.ID)

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			monthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, month, category.ID)
			require.Nil(t, err)

			assert.Equal(t,
				monthCategory,
				beans.MonthCategory{
					ID:         monthCategory.ID,
					CategoryID: category.ID,
					MonthID:    month.ID,
					Amount:     beans.NewAmount(5, 0),
				},
			)
		})

		t.Run("uses existing month category", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			auth := testutils.BudgetAuthContext(t, userID, budget)

			group := factory.MakeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Category", group.ID, budget.ID)
			monthCategory := factory.MakeMonthCategory(month.ID, category.ID, beans.NewAmount(4, 0))

			err := c.SetCategoryAmount(context.Background(), auth, month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			dbMonthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, month, category.ID)
			require.Nil(t, err)

			assert.Equal(t,
				dbMonthCategory,
				beans.MonthCategory{
					ID:         monthCategory.ID,
					CategoryID: category.ID,
					MonthID:    month.ID,
					Amount:     beans.NewAmount(5, 0),
				},
			)
		})

	})
}
