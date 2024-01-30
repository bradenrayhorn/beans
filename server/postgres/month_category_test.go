package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonthCategory(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()
	ds := postgres.NewDataSource(pool)
	factory := testutils.Factory(t, ds)

	txManager := postgres.NewTxManager(pool)
	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)

	userID := factory.MakeUser("user")
	budgetID := factory.MakeBudget("budget", userID).ID
	budgetID2 := factory.MakeBudget("budget2", userID).ID
	account := factory.MakeAccount("account", budgetID)
	groupID := factory.MakeCategoryGroup("group", budgetID).ID
	categoryID := factory.MakeCategory("group", groupID, budgetID).ID
	categoryID2 := factory.MakeCategory("group", groupID, budgetID).ID
	monthMarch := factory.MakeMonth(budgetID, testutils.NewDate(t, "2022-03-01"))
	monthApril := factory.MakeMonth(budgetID, testutils.NewDate(t, "2022-04-01"))
	monthMay := factory.MakeMonth(budgetID, testutils.NewDate(t, "2022-05-01"))
	monthJune := factory.MakeMonth(budgetID, testutils.NewDate(t, "2022-06-01"))

	budget2Month := factory.MakeMonth(budgetID2, testutils.NewDate(t, "2022-05-01"))
	budget2GroupID := factory.MakeCategoryGroup("group", budgetID2).ID
	budget2CategoryID := factory.MakeCategory("group", budget2GroupID, budgetID2).ID
	budget2Account := factory.MakeAccount("account", budgetID2)

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate month_categories cascade; truncate transactions cascade;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthMay.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(1, 0),

			// Values only returned on get
			Activity:  beans.NewAmount(0, 0),
			Available: beans.NewAmount(1, 0),
		}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("can create with empty amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewEmptyAmount()}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, monthCategory.ID, res[0].ID)
		assert.Equal(t, beans.NewAmount(0, 0), res[0].Amount)
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))
		assertPgError(t, pgerrcode.UniqueViolation, monthCategoryRepository.Create(context.Background(), nil, monthCategory))
	})

	t.Run("cannot create with duplicate month and category", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))
		monthCategory.ID = beans.NewBeansID()
		assertPgError(t, pgerrcode.UniqueViolation, monthCategoryRepository.Create(context.Background(), nil, monthCategory))
	})

	t.Run("create respects tx", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}

		// make transaction
		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// create but do not commit
		require.Nil(t, monthCategoryRepository.Create(context.Background(), tx, monthCategory))

		// try to find category, should fail
		categories, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, categories, 0)

		// commit
		require.Nil(t, tx.Commit(context.Background()))

		// try to find month, should succeed
		categories, err = monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, categories, 1)
		require.Equal(t, monthCategory.ID, categories[0].ID)
	})

	t.Run("can update amount", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthMay.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(1, 0),

			// Values only returned on get
			Activity:  beans.NewAmount(0, 0),
			Available: beans.NewAmount(1, 0),
		}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))

		require.Nil(t, monthCategoryRepository.UpdateAmount(context.Background(), monthCategory.ID, beans.NewAmount(5, -1)))
		monthCategory.Amount = beans.NewAmount(5, -1)
		monthCategory.Available = beans.NewAmount(5, -1)

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory, res[0]))
	})

	t.Run("get filters by month", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthMay.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(1, 0),

			// Values only returned on get
			Activity:  beans.NewAmount(0, 0),
			Available: beans.NewAmount(1, 0),
		}
		monthCategory2 := &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthJune.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(1, 0),

			// Values only returned on get
			Activity:  beans.NewAmount(0, 0),
			Available: beans.NewAmount(1, 0),
		}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory2))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.True(t, reflect.DeepEqual(monthCategory1, res[0]))
	})

	t.Run("sums activity and available properly", func(t *testing.T) {
		defer cleanup()
		monthCategoryMarch1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMarch.ID, CategoryID: categoryID, Amount: beans.NewAmount(5, 0)}
		monthCategoryApril1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthApril.ID, CategoryID: categoryID, Amount: beans.NewAmount(15, 0)}
		monthCategoryMay1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategoryJune1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID, Amount: beans.NewAmount(15, 0)}

		monthCategoryMay2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID2, Amount: beans.NewAmount(5, 0)}
		month2CategoryMay1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: budget2Month.ID, CategoryID: budget2CategoryID, Amount: beans.NewAmount(1, 0)}

		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategoryMay1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategoryMarch1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategoryApril1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategoryJune1))

		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategoryMay2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, month2CategoryMay1))

		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-5, 0),
			Date:       testutils.NewDate(t, "2022-03-20"),
		})
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(-6, 0),
			Date:       testutils.NewDate(t, "2022-03-20"),
		})

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
		makeTransaction(t, pool, &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  budget2Account.ID,
			CategoryID: budget2CategoryID,
			Amount:     beans.NewAmount(-4, 0),
			Date:       testutils.NewDate(t, "2022-05-01"),
		})

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 2)

		category1 := findResult(res, func(m *beans.MonthCategory) bool { return m.CategoryID == categoryID })
		category2 := findResult(res, func(m *beans.MonthCategory) bool { return m.CategoryID == categoryID2 })

		// category 1
		require.Equal(t, beans.NewAmount(-1182, -2), category1.Activity) // -5 - 6.82 (2 May transactions)
		require.Equal(t, beans.NewAmount(1, 0), category1.Amount)        // 1 (Amount assigned in May)
		require.Equal(t, beans.NewAmount(-682, -2), category1.Available) // (5 + 15 + 1) (Assigned) (- 5 - 6) (March Spending) (- 5 - 6.82) (May Spending) - 5 (April Spending)

		// category 2
		require.Equal(t, beans.NewAmount(5, 0), category2.Amount)
		require.Equal(t, beans.NewAmount(0, 0), category2.Activity)
		require.Equal(t, beans.NewAmount(5, 0), category2.Available)
	})

	t.Run("sums spent no transactions to zero", func(t *testing.T) {
		defer cleanup()
		monthCategory := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory))

		res, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.Equal(t, beans.NewAmount(0, 0), res[0].Activity)
	})

	t.Run("get or create respects month and category", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory4))

		res, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, monthMay.ID, categoryID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(monthCategory1, res))
	})

	t.Run("get or create returns new", func(t *testing.T) {
		defer cleanup()
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID, Amount: beans.NewAmount(1, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID2, Amount: beans.NewAmount(1, 0)}
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory4))

		existingIDs := []beans.ID{monthCategory2.ID, monthCategory3.ID, monthCategory4.ID}

		monthCategory, err := monthCategoryRepository.GetOrCreate(context.Background(), nil, monthMay.ID, categoryID)
		require.Nil(t, err)

		assert.NotContains(t, existingIDs, monthCategory.ID)
		assert.True(t, reflect.DeepEqual(
			monthCategory,
			&beans.MonthCategory{
				ID:         monthCategory.ID,
				MonthID:    monthMay.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			},
		))
	})

	t.Run("get or create respects tx", func(t *testing.T) {
		defer cleanup()

		// make transaction
		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// get or create but do not commit
		_, err = monthCategoryRepository.GetOrCreate(context.Background(), tx, monthMay.ID, categoryID)
		require.Nil(t, err)

		// try to find category, should fail
		categories, err := monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, categories, 0)

		// commit
		require.Nil(t, tx.Commit(context.Background()))

		// try to find month, should succeed
		categories, err = monthCategoryRepository.GetForMonth(context.Background(), monthMay)
		require.Nil(t, err)
		require.Len(t, categories, 1)
	})

	t.Run("can get assigned in month", func(t *testing.T) {
		defer cleanup()
		monthCategory1 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthMay.ID, CategoryID: categoryID, Amount: beans.NewAmount(7, 0)}
		monthCategory2 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID2, Amount: beans.NewAmount(8, 0)}
		monthCategory3 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: monthJune.ID, CategoryID: categoryID, Amount: beans.NewAmount(3, 0)}

		monthCategory4 := &beans.MonthCategory{ID: beans.NewBeansID(), MonthID: budget2Month.ID, CategoryID: budget2CategoryID, Amount: beans.NewAmount(9, 0)}

		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory1))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory2))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory3))
		require.Nil(t, monthCategoryRepository.Create(context.Background(), nil, monthCategory4))

		amount, err := monthCategoryRepository.GetAssignedInMonth(context.Background(), monthJune.ID)
		require.Nil(t, err)
		assert.Equal(t, beans.NewAmount(11, 0), amount)
	})

	t.Run("can get blank assigned in month", func(t *testing.T) {
		defer cleanup()

		amount, err := monthCategoryRepository.GetAssignedInMonth(context.Background(), monthJune.ID)
		require.Nil(t, err)
		assert.Equal(t, beans.NewAmount(0, 0), amount)
	})
}
