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
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID
	budgetID2 := testutils.MakeBudget(t, pool, "budget2", userID).ID
	account := testutils.MakeAccount(t, pool, "account", budgetID)
	groupID := testutils.MakeCategoryGroup(t, pool, "group", budgetID).ID
	categoryID := testutils.MakeCategory(t, pool, "group", groupID, budgetID).ID
	categoryID2 := testutils.MakeCategory(t, pool, "group", groupID, budgetID).ID
	month := testutils.MakeMonth(t, pool, budgetID, testutils.NewDate(t, "2022-05-01"))
	month2 := testutils.MakeMonth(t, pool, budgetID, testutils.NewDate(t, "2022-06-01"))

	budget2Month := testutils.MakeMonth(t, pool, budgetID2, testutils.NewDate(t, "2022-05-01"))
	budget2GroupID := testutils.MakeCategoryGroup(t, pool, "group", budgetID2).ID
	budget2CategoryID := testutils.MakeCategory(t, pool, "group", budget2GroupID, budgetID2).ID

	cleanup := func() {
		pool.Exec(context.Background(), "truncate month_categories cascade; truncate transactions cascade;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0), Activity: beans.NewAmount(0, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("can create with empty amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewEmptyAmount()}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, monthCategory.ID, res[0].ID)
		assert.True(t, res[0].Amount.Empty())
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))
		assertPgError(t, pgerrcode.UniqueViolation, monthCategoryRepository.Create(context.Background(), monthCategory))
	})

	t.Run("can update amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0), Activity: beans.NewAmount(0, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		require.Nil(t, monthCategoryRepository.UpdateAmount(context.Background(), monthCategory.ID, beans.NewAmount(5, -1)))
		monthCategory.Amount = beans.NewAmount(5, -1)

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("get filters by month", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0), Activity: beans.NewAmount(0, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0), Activity: beans.NewAmount(0, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory1, res[0]))
	})

	t.Run("sums multiple transactions in activity", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-5, 0),
			Date:       testutils.NewDate(t, "2022-05-20"),
		})
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-682, -2),
			Date:       testutils.NewDate(t, "2022-05-20"),
		})

		// make transaction in past and future months that are not included in total
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-5, 0),
			Date:       testutils.NewDate(t, "2022-04-30"),
		})
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-5, 0),
			Date:       testutils.NewDate(t, "2022-06-01"),
		})

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Equal(t, beans.NewAmount(-1182, -2), res[0].Activity)
	})

	t.Run("sums spent no transactions to zero", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Equal(t, beans.NewAmount(0, 0), res[0].Activity)
	})

	t.Run("get or create respects month and category", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory4))

		res, err := monthCategoryRepository.GetOrCreate(context.Background(), month.ID, categoryID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(monthCategory1, res))
	})

	t.Run("get or create returns new", func(t *testing.T) {
		defer cleanup()
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory4))

		existingIDs := []beans.ID{monthCategory2.ID, monthCategory3.ID, monthCategory4.ID}

		monthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), month.ID, categoryID)
		require.Nil(t, err)

		assert.NotContains(t, existingIDs, monthCategory.ID)
		assert.True(t, reflect.DeepEqual(
			monthCategory,
			&beans.MonthCategory{
				ID:         monthCategory.ID,
				MonthID:    month.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			},
		))
	})

	t.Run("can get amount in budget", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month.ID, CategoryID: categoryID, Amount: beans.NewAmount(7, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: month2.ID, CategoryID: categoryID2, Amount: beans.NewAmount(8, 0)}

		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: budget2Month.ID, CategoryID: budget2CategoryID, Amount: beans.NewAmount(9, 0)}

		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), monthCategory3))

		amount, err := monthCategoryRepository.GetAmountInBudget(context.Background(), budgetID)
		require.Nil(t, err)

		assert.Equal(t, beans.NewAmount(15, 0), amount)
	})
}
