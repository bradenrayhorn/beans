package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonthCategory(t *testing.T) {
	pool, stop := testutils.StartPool(t)
	defer stop()

	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)
	groupID := makeCategoryGroup(t, pool, "group", budgetID)
	categoryID := makeCategory(t, pool, "group", groupID, budgetID)
	categoryID2 := makeCategory(t, pool, "group", groupID, budgetID)
	monthID := makeMonth(t, pool, budgetID, testutils.NewDate(t, "2022-05-01"))
	monthID2 := makeMonth(t, pool, budgetID, testutils.NewDate(t, "2022-06-01"))

	cleanup := func() {
		pool.Exec(context.Background(), "truncate month_categories cascade;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("can create with empty amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewEmptyAmount()}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, monthCategory.ID, res[0].ID)
		assert.True(t, res[0].Amount.Empty())
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))
		assertPgError(t, pgerrcode.UniqueViolation, monthCategoryRepository.Create(context.Background(), monthCategory))
	})

	t.Run("can update amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		require.Nil(t, monthCategoryRepository.UpdateAmount(context.Background(), monthCategory.ID, beans.NewAmount(5, -1)))
		monthCategory.Amount = beans.NewAmount(5, -1)

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("get filters by month", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID2, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory1, res[0]))
	})

	t.Run("can get by month and category", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID2, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthID2, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory4))

		res, err := monthCategoryRepository.GetByMonthAndCategory(context.Background(), monthID, categoryID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(monthCategory1, res))
	})

	t.Run("get by month and category returns not found error", func(t *testing.T) {
		defer cleanup()

		_, err := monthCategoryRepository.GetByMonthAndCategory(context.Background(), monthID, categoryID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})
}
