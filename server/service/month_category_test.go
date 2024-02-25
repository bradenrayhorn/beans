package service_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonthCategory(t *testing.T) {
	services, factory, _, _ := makeServices(t)
	ctx := context.Background()

	t.Run("GetForMonth", func(t *testing.T) {

		t.Run("with a new month", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			month := factory.Month(beans.Month{BudgetID: budget.ID})

			res, err := services.MonthCategory.GetForMonth(ctx, month)
			require.NoError(t, err)
			assert.Equal(t, 0, len(res))
		})

		t.Run("if everything is zero", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			month := factory.Month(beans.Month{BudgetID: budget.ID})
			category := factory.Category(beans.Category{})
			monthCategory := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: month.ID, CategoryID: category.ID})

			res, err := services.MonthCategory.GetForMonth(ctx, month)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, beans.MonthCategoryWithDetails{
				ID:         monthCategory.ID,
				CategoryID: category.ID,
				Amount:     beans.NewAmount(0, 0),
				Activity:   beans.NewAmount(0, 0),
				Available:  beans.NewAmount(0, 0),
			}, res[0])
		})

		t.Run("includes previous assigned in available", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			april := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-04-01")})
			may := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-05-01")})
			june := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-06-01")})
			category := factory.Category(beans.Category{})

			// assign $3 in April, $0 in May, $5 in June
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category.ID, Amount: beans.NewAmount(3, 0)})
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: may.ID, CategoryID: category.ID})
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: june.ID, CategoryID: category.ID, Amount: beans.NewAmount(5, 0)})

			// get for May
			res, err := services.MonthCategory.GetForMonth(ctx, may)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))

			// there should be $3 available in May
			assert.Equal(t, beans.NewAmount(0, 0), res[0].Amount)
			assert.Equal(t, beans.NewAmount(3, 0), res[0].Available)
			assert.Equal(t, beans.NewAmount(0, 0), res[0].Activity)
		})

		t.Run("includes activity in available", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			april := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-04-01")})
			category := factory.Category(beans.Category{})
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category.ID})

			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-03-31"),
				Amount:     beans.NewAmount(-3, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-04-01"),
				Amount:     beans.NewAmount(-5, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-04-30"),
				Amount:     beans.NewAmount(-6, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-05-01"),
				Amount:     beans.NewAmount(-9, 0),
			})

			// get for April
			res, err := services.MonthCategory.GetForMonth(ctx, april)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))

			// there should be -$14 available, -$11 in April
			assert.Equal(t, beans.NewAmount(0, 0), res[0].Amount)
			assert.Equal(t, beans.NewAmount(-14, 0), res[0].Available)
			assert.Equal(t, beans.NewAmount(-11, 0), res[0].Activity)
		})

		t.Run("includes current assigned in available", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			may := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-05-01")})
			category := factory.Category(beans.Category{})

			// assign $3 in May
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: may.ID, CategoryID: category.ID, Amount: beans.NewAmount(3, 0)})

			// get for May
			res, err := services.MonthCategory.GetForMonth(ctx, may)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))

			// there should be $3 available and assigned in May
			assert.Equal(t, beans.NewAmount(3, 0), res[0].Amount)
			assert.Equal(t, beans.NewAmount(3, 0), res[0].Available)
			assert.Equal(t, beans.NewAmount(0, 0), res[0].Activity)
		})

		t.Run("sums and groups", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			march := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-03-01")})
			april := factory.Month(beans.Month{BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-04-01")})
			category1 := factory.Category(beans.Category{})
			category2 := factory.Category(beans.Category{})

			// assign $5 to category1 for March, none to category2
			factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: march.ID, CategoryID: category1.ID, Amount: beans.NewAmount(5, 0)})

			// assign $3 to category1 for April, $2 to category2
			aprilCategory1 := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category1.ID, Amount: beans.NewAmount(3, 0)})
			aprilCategory2 := factory.MonthCategory(budget.ID, beans.MonthCategory{MonthID: april.ID, CategoryID: category2.ID, Amount: beans.NewAmount(2, 0)})

			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category1.ID, // spend $3 in March category1
				Date:       testutils.NewDate(t, "2022-03-31"),
				Amount:     beans.NewAmount(-3, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category2.ID, // spend $1 in March category2
				Date:       testutils.NewDate(t, "2022-03-31"),
				Amount:     beans.NewAmount(-1, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category1.ID, // spend $5 in April category1
				Date:       testutils.NewDate(t, "2022-04-01"),
				Amount:     beans.NewAmount(-5, 0),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				CategoryID: category2.ID, // spend $3 in April category2
				Date:       testutils.NewDate(t, "2022-04-30"),
				Amount:     beans.NewAmount(-3, 0),
			})

			// get for April and verify
			res, err := services.MonthCategory.GetForMonth(ctx, april)
			require.NoError(t, err)
			assert.Equal(t, 2, len(res))
			assert.ElementsMatch(t, []beans.MonthCategoryWithDetails{
				{
					ID:         aprilCategory1.ID,
					CategoryID: category1.ID,
					Amount:     beans.NewAmount(3, 0),
					Available:  beans.NewAmount(0, 0),
					Activity:   beans.NewAmount(-5, 0),
				},
				{
					ID:         aprilCategory2.ID,
					CategoryID: category2.ID,
					Amount:     beans.NewAmount(2, 0),
					Available:  beans.NewAmount(-2, 0),
					Activity:   beans.NewAmount(-3, 0),
				},
			}, res)
		})
	})

}
