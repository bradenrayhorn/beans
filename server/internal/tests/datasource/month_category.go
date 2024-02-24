package datasource

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMonthCategory(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)

	monthCategoryRepository := ds.MonthCategoryRepository()
	ctx := context.Background()

	t.Run("can create and get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		monthCategory := beans.MonthCategory{
			ID:         beans.NewID(),
			MonthID:    month.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(1, 0),
		}
		require.Nil(t, monthCategoryRepository.Create(ctx, nil, monthCategory))

		res, err := monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, res, 1)

		assert.Equal(t, beans.MonthCategory{
			ID:         monthCategory.ID,
			MonthID:    month.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(1, 0),
		}, res[0])
	})

	t.Run("can create with empty amount", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		monthCategory := beans.MonthCategory{
			ID:         beans.NewID(),
			MonthID:    month.ID,
			CategoryID: category.ID,
			Amount:     beans.NewEmptyAmount(),
		}
		require.Nil(t, monthCategoryRepository.Create(ctx, nil, monthCategory))

		res, err := monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, monthCategory.ID, res[0].ID)
		assert.Equal(t, beans.NewAmount(0, 0), res[0].Amount)
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		monthCategory := factory.MonthCategory(budget.ID, beans.MonthCategory{})

		assert.NotNil(t, monthCategoryRepository.Create(ctx, nil, monthCategory))
	})

	t.Run("cannot create with duplicate month and category", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		monthCategory := factory.MonthCategory(budget.ID, beans.MonthCategory{})

		// try with a new month category id
		monthCategory.ID = beans.NewID()

		assert.NotNil(t, monthCategoryRepository.Create(ctx, nil, monthCategory))
	})

	t.Run("create respects tx", func(t *testing.T) {
		txManager := ds.TxManager()

		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		monthCategory := beans.MonthCategory{
			ID:         beans.NewID(),
			MonthID:    month.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(1, 0),
		}

		// make transaction
		tx, err := txManager.Create(ctx)
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// create but do not commit
		require.Nil(t, monthCategoryRepository.Create(ctx, tx, monthCategory))

		// try to find category, should fail
		categories, err := monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, categories, 0)

		// commit
		require.Nil(t, tx.Commit(ctx))

		// try to find month, should succeed
		categories, err = monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, categories, 1)
		require.Equal(t, monthCategory.ID, categories[0].ID)
	})

	t.Run("can update amount", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})
		monthCategory := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month.ID})

		// update amount
		monthCategory.Amount = beans.NewAmount(5, -1)
		require.Nil(t, monthCategoryRepository.UpdateAmount(ctx, monthCategory))

		// get and verify amount is updated
		res, err := monthCategoryRepository.GetOrCreate(ctx, nil, month, monthCategory.CategoryID)
		require.Nil(t, err)

		assert.Equal(t, beans.NewAmount(5, -1), res.Amount)
	})

	t.Run("get for month", func(t *testing.T) {

		t.Run("can get for month", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			month := factory.Month(beans.Month{BudgetID: budget.ID})

			monthCategory := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month.ID, Amount: beans.NewAmount(1, 0)})

			// get and verify results
			res, err := monthCategoryRepository.GetForMonth(ctx, month)
			require.NoError(t, err)

			assert.Equal(t, 1, len(res))
			assert.Equal(t, beans.MonthCategory{
				ID:         monthCategory.ID,
				MonthID:    month.ID,
				CategoryID: monthCategory.CategoryID,
				Amount:     beans.NewAmount(1, 0),
			}, res[0])
		})

		t.Run("get for month filters by month", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			month1 := factory.Month(beans.Month{BudgetID: budget.ID})
			month2 := factory.Month(beans.Month{BudgetID: budget.ID})

			// this month category should not be returned
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month2.ID})

			res, err := monthCategoryRepository.GetForMonth(ctx, month1)
			require.NoError(t, err)
			require.Equal(t, 0, len(res))
		})
	})

	t.Run("get assigned by category", func(t *testing.T) {

		t.Run("sums and groups", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			category1 := factory.Category(beans.Category{})
			category2 := factory.Category(beans.Category{})

			april := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-04-01")})
			march := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-03-01")})

			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category1.ID, Amount: beans.NewAmount(2, 0)}) // $2 in april C1
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: march.ID, CategoryID: category1.ID, Amount: beans.NewAmount(3, 0)}) // $3 in march C1
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category2.ID, Amount: beans.NewAmount(6, 0)}) // $6 in april C2

			res, err := monthCategoryRepository.GetAssignedByCategory(ctx, budget.ID, testutils.NewDate(t, "2022-05-01"))
			require.NoError(t, err)

			assert.Equal(t, 2, len(res))
			assert.Equal(t, beans.NewAmount(5, 0), res[category1.ID])
			assert.Equal(t, beans.NewAmount(6, 0), res[category2.ID])
		})

		t.Run("filters by date", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			month := factory.Month(beans.Month{Date: testutils.NewMonthDate(t, "2022-05-01")})
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month.ID, Amount: beans.NewAmount(1, 0)})

			// the month category is in May, we are filtering to before May
			before := testutils.NewDate(t, "2022-05-01")
			res, err := monthCategoryRepository.GetAssignedByCategory(ctx, budget.ID, before)
			require.NoError(t, err)
			assert.Equal(t, 0, len(res))
		})

		t.Run("filters by budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			month := factory.Month(beans.Month{Date: testutils.NewMonthDate(t, "2022-04-01")})
			factory.MonthCategory(budget2.ID, beans.MonthCategory{MonthID: month.ID})

			before := testutils.NewDate(t, "2022-05-01")
			res, err := monthCategoryRepository.GetAssignedByCategory(ctx, budget.ID, before)
			require.NoError(t, err)
			assert.Equal(t, 0, len(res))
		})
	})

	t.Run("get or create respects month and category", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month1 := factory.Month(beans.Month{BudgetID: budget.ID})
		month2 := factory.Month(beans.Month{BudgetID: budget.ID})

		category1 := factory.Category(beans.Category{BudgetID: budget.ID})
		category2 := factory.Category(beans.Category{BudgetID: budget.ID})

		expected := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month1.ID, CategoryID: category1.ID})
		factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month2.ID, CategoryID: category1.ID})
		factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month1.ID, CategoryID: category2.ID})
		factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month2.ID, CategoryID: category2.ID})

		res, err := monthCategoryRepository.GetOrCreate(ctx, nil, month1, category1.ID)
		require.Nil(t, err)
		assert.Equal(t, expected.ID, res.ID)
	})

	t.Run("get or create returns new", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month1 := factory.Month(beans.Month{BudgetID: budget.ID})
		month2 := factory.Month(beans.Month{BudgetID: budget.ID})

		category1 := factory.Category(beans.Category{BudgetID: budget.ID})
		category2 := factory.Category(beans.Category{BudgetID: budget.ID})

		mc1 := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month2.ID, CategoryID: category1.ID})
		mc2 := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month1.ID, CategoryID: category2.ID})
		mc3 := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month2.ID, CategoryID: category2.ID})

		existingIDs := []beans.ID{mc1.ID, mc2.ID, mc3.ID}

		monthCategory, err := monthCategoryRepository.GetOrCreate(ctx, nil, month1, category1.ID)
		require.Nil(t, err)

		assert.NotContains(t, existingIDs, monthCategory.ID)
		assert.Equal(t,
			monthCategory,
			beans.MonthCategory{
				ID:         monthCategory.ID,
				MonthID:    month1.ID,
				CategoryID: category1.ID,
				Amount:     beans.NewAmount(0, 0),
			},
		)
	})

	t.Run("get or create respects tx", func(t *testing.T) {
		txManager := ds.TxManager()

		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		// make transaction
		tx, err := txManager.Create(ctx)
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// get or create but do not commit
		_, err = monthCategoryRepository.GetOrCreate(ctx, tx, month, category.ID)
		require.Nil(t, err)

		// try to find, should fail
		categories, err := monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, categories, 0)

		// commit
		require.Nil(t, tx.Commit(ctx))

		// try to find, should succeed
		categories, err = monthCategoryRepository.GetForMonth(ctx, month)
		require.Nil(t, err)
		require.Len(t, categories, 1)
	})

	t.Run("can get assigned in month", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		month1 := factory.Month(beans.Month{BudgetID: budget1.ID})
		month2 := factory.Month(beans.Month{BudgetID: budget1.ID})

		factory.MonthCategory(budget1.ID, beans.MonthCategory{MonthID: month1.ID, Amount: beans.NewAmount(7, 0)})
		factory.MonthCategory(budget1.ID, beans.MonthCategory{MonthID: month2.ID, Amount: beans.NewAmount(8, 0)})
		factory.MonthCategory(budget1.ID, beans.MonthCategory{MonthID: month2.ID, Amount: beans.NewAmount(3, 0)})

		factory.MonthCategory(budget2.ID, beans.MonthCategory{Amount: beans.NewAmount(9, 0)})

		amount, err := monthCategoryRepository.GetAssignedInMonth(ctx, month2)
		require.Nil(t, err)
		assert.Equal(t, beans.NewAmount(11, 0), amount)
	})

	t.Run("can get blank assigned in month", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})

		amount, err := monthCategoryRepository.GetAssignedInMonth(ctx, month)
		require.Nil(t, err)
		assert.Equal(t, beans.NewAmount(0, 0), amount)
	})
}
