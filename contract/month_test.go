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

	monthRepository := postgres.NewMonthRepository(pool)
	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)
	c := contract.NewMonthContract(monthRepository, monthCategoryRepository)

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get non existant month", func(t *testing.T) {
			defer cleanup()

			_, _, err := c.Get(context.Background(), beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("can get month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)
			monthCategory := testutils.MakeMonthCategory(t, pool, month.ID, category.ID, beans.NewAmount(34, -1))

			dbMonth, dbCategories, err := c.Get(context.Background(), month.ID)
			require.Nil(t, err)

			assert.True(t, reflect.DeepEqual(month, dbMonth))
			require.Len(t, dbCategories, 1)

			monthCategory.Spent = beans.NewAmount(0, 0)
			assert.True(t, reflect.DeepEqual(monthCategory, dbCategories[0]))
		})
	})

	t.Run("create", func(t *testing.T) {
		t.Run("creates new month", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)

			date := testutils.NewMonthDate(t, "2022-05-01")

			month, err := c.CreateMonth(context.Background(), budget.ID, date)
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
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			returnedMonth, err := c.CreateMonth(context.Background(), budget.ID, month.Date)
			require.Nil(t, err)

			// month was returned
			assert.True(t, reflect.DeepEqual(month, returnedMonth))
		})
	})

	t.Run("set category amount", func(t *testing.T) {
		t.Run("amount must be not be zero", func(t *testing.T) {
			defer cleanup()

			err := c.SetCategoryAmount(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(0, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("amount must be not be negative", func(t *testing.T) {
			defer cleanup()

			err := c.SetCategoryAmount(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(-5, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("creates new month category", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)

			err := c.SetCategoryAmount(context.Background(), month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			monthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), month.ID, category.ID)
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
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)
			monthCategory := testutils.MakeMonthCategory(t, pool, month.ID, category.ID, beans.NewAmount(4, 0))

			err := c.SetCategoryAmount(context.Background(), month.ID, category.ID, beans.NewAmount(5, 0))
			require.Nil(t, err)

			dbMonthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), month.ID, category.ID)
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