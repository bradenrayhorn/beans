package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMonth(t *testing.T, interactor Interactor) {

	t.Run("get or create", func(t *testing.T) {
		t.Run("can create a new month", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create a new month
			date := testutils.NewMonthDate(t, "2024-01-01")
			month, err := interactor.MonthGetOrCreate(t, c.ctx, date)
			require.NoError(t, err)

			// check month details
			assert.NotEmpty(t, month.ID)
			assert.Equal(t, date, month.Date)
			assert.Equal(t, beans.NewAmount(0, 0), month.Carryover)
			assert.Equal(t, beans.NewAmount(0, 0), month.CarriedOver)
			assert.Equal(t, beans.NewAmount(0, 0), month.Income)
			assert.Equal(t, beans.NewAmount(0, 0), month.Assigned)
			assert.Equal(t, beans.NewAmount(0, 0), month.Budgetable)

			// only the income category should exist
			assert.Equal(t, 1, len(month.Categories))

			// check the income category
			monthCategory := month.Categories[0]
			assert.NotEmpty(t, monthCategory.ID)
			assert.Equal(t, c.findIncomeCategory().ID, monthCategory.CategoryID)
			assert.Equal(t, beans.NewAmount(0, 0), monthCategory.Amount)
			assert.Equal(t, beans.NewAmount(0, 0), monthCategory.Activity)
			assert.Equal(t, beans.NewAmount(0, 0), monthCategory.Available)
		})

		t.Run("reuses an existing month", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create a new month
			date := testutils.NewMonthDate(t, "2024-01-01")
			month, err := interactor.MonthGetOrCreate(t, c.ctx, date)
			require.NoError(t, err)

			// try making the same month again
			res, err := interactor.MonthGetOrCreate(t, c.ctx, date)
			require.NoError(t, err)

			// the result should have the same ID
			assert.Equal(t, month.ID, res.ID)
		})

		t.Run("does not reuse a month from another budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			// create a new month in c1
			date := testutils.NewMonthDate(t, "2024-01-01")
			month, err := interactor.MonthGetOrCreate(t, c1.ctx, date)
			require.NoError(t, err)

			// try making the same month in c2
			res, err := interactor.MonthGetOrCreate(t, c2.ctx, date)
			require.NoError(t, err)

			// the month ID should be different
			assert.NotEqual(t, month.ID, res.ID)
		})

		t.Run("can get a month with many details", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create May and April with carryovers
			monthApril := c.Month(MonthOpts{Date: "2022-04-01", Carryover: "6.7"})
			monthMay := c.Month(MonthOpts{Date: "2022-05-01", Carryover: "0.4"})

			// create account and categories
			account := c.Account(AccountOpts{})
			categoryIncome := c.findIncomeCategory()
			categoryBills := c.Category(CategoryOpts{})

			// assign to bills category
			c.setAssigned(monthMay, categoryBills, "3.4")
			c.setAssigned(monthApril, categoryBills, "3.4")

			// create transactions
			c.Transaction(TransactionOpts{ // earned $6 in March
				Account:  account,
				Category: categoryIncome,
				Date:     "2022-03-01",
				Amount:   "6",
			})
			c.Transaction(TransactionOpts{ // earned $9 in May
				Account:  account,
				Category: categoryIncome,
				Date:     "2022-05-01",
				Amount:   "9",
			})
			c.Transaction(TransactionOpts{ // earned $9 in June
				Account:  account,
				Category: categoryIncome,
				Date:     "2022-06-01",
				Amount:   "3",
			})

			// get month
			res, err := interactor.MonthGetOrCreate(t, c.ctx, monthMay.Date)
			require.NoError(t, err)

			// check month details
			assert.Equal(t, monthMay.ID, res.ID)
			assert.Equal(t, monthMay.Date, res.Date)
			assert.Equal(t, beans.NewAmount(4, -1), res.Carryover)    // carrying over $0.4
			assert.Equal(t, beans.NewAmount(67, -1), res.CarriedOver) // carried over $6.7 from April
			assert.Equal(t, beans.NewAmount(9, 0), res.Income)        // earned $9 in May
			assert.Equal(t, beans.NewAmount(34, -1), res.Assigned)    // assigned $3.4 to bills category in May
			assert.Equal(t, beans.NewAmount(119, -1), res.Budgetable) // $6.7 (april carryover) + $9 (earnings) - $0.4 (carrying over) - $3.4 (assigned) = $11.9

			// the income and bills should exist
			assert.Equal(t, 2, len(res.Categories))

			// check the income category
			findMonthCategory(t, res.Categories, categoryIncome.ID, func(it beans.MonthCategoryWithDetails) {
				assert.NotEmpty(t, it.ID)
				assert.Equal(t, beans.NewAmount(0, 0), it.Amount)
				assert.Equal(t, beans.NewAmount(9, 0), it.Activity)   // earned $9
				assert.Equal(t, beans.NewAmount(15, 0), it.Available) // May + April income
			})

			// check bills category
			findMonthCategory(t, res.Categories, categoryBills.ID, func(it beans.MonthCategoryWithDetails) {
				assert.NotEmpty(t, it.ID)
				assert.Equal(t, beans.NewAmount(34, -1), it.Amount)    // Assigned $3.4
				assert.Equal(t, beans.NewAmount(0, 0), it.Activity)    // Spent $0
				assert.Equal(t, beans.NewAmount(68, -1), it.Available) // Have $6.8 (assigned + assigned in April)
			})
		})
	})

	t.Run("update", func(t *testing.T) {

		t.Run("cannot update a month that does not exist", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			err := interactor.MonthUpdate(t, c.ctx, beans.NewID(), beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot update month from another budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			month := c2.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthUpdate(t, c1.ctx, month.ID, beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot carryover a negative amount", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			month := c.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthUpdate(t, c.ctx, month.ID, beans.NewAmount(-1, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot carryover nothing", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			month := c.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthUpdate(t, c.ctx, month.ID, beans.NewEmptyAmount())
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can carryover an amount", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			month := c.Month(MonthOpts{Date: "2022-04-01"})

			// carryover $1
			err := interactor.MonthUpdate(t, c.ctx, month.ID, beans.NewAmount(1, 0))
			require.NoError(t, err)

			// get month and check carryover
			res, err := interactor.MonthGetOrCreate(t, c.ctx, testutils.NewMonthDate(t, "2022-04-01"))
			require.NoError(t, err)
			assert.Equal(t, beans.NewAmount(1, 0), res.Carryover)
		})

		t.Run("can carryover zero", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			month := c.Month(MonthOpts{Date: "2022-04-01", Carryover: "10.5"})

			// carryover $1
			err := interactor.MonthUpdate(t, c.ctx, month.ID, beans.NewAmount(0, 0))
			require.NoError(t, err)

			// get month and check carryover
			res, err := interactor.MonthGetOrCreate(t, c.ctx, testutils.NewMonthDate(t, "2022-04-01"))
			require.NoError(t, err)
			assert.Equal(t, beans.NewAmount(0, 0), res.Carryover)
		})
	})

	t.Run("set category amount", func(t *testing.T) {

		t.Run("cannot assign with a month that does not exist", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			category := c.Category(CategoryOpts{})

			err := interactor.MonthSetCategoryAmount(t, c.ctx, beans.NewID(), category.ID, beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot assign with month from another budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			category := c1.Category(CategoryOpts{})
			month := c2.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthSetCategoryAmount(t, c1.ctx, month.ID, category.ID, beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot assign with category that does not exist", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			month := c.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, beans.NewID(), beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot assign with category from another budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			category := c2.Category(CategoryOpts{})
			month := c1.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthSetCategoryAmount(t, c1.ctx, month.ID, category.ID, beans.NewAmount(1, 0))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot assign a negative amount", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			category := c.Category(CategoryOpts{})
			month := c.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, category.ID, beans.NewAmount(-1, 0))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot assign nothing", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			category := c.Category(CategoryOpts{})
			month := c.Month(MonthOpts{Date: "2022-04-01"})

			err := interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, category.ID, beans.NewEmptyAmount())
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can assign an amount", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			category := c.Category(CategoryOpts{})
			month := c.Month(MonthOpts{Date: "2022-04-01"})

			// assign $1
			err := interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, category.ID, beans.NewAmount(1, 0))
			require.NoError(t, err)

			// check month and see if month category is assigned
			res, err := interactor.MonthGetOrCreate(t, c.ctx, testutils.NewMonthDate(t, "2022-04-01"))
			require.NoError(t, err)

			assert.Equal(t, 2, len(res.Categories)) // income + expense category
			findMonthCategory(t, res.Categories, category.ID, func(it beans.MonthCategoryWithDetails) {
				assert.Equal(t, beans.NewAmount(1, 0), it.Amount)
			})
		})

		t.Run("can update a month category assigned", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			category := c.Category(CategoryOpts{})
			month := c.Month(MonthOpts{Date: "2022-04-01"})

			// assign $1
			err := interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, category.ID, beans.NewAmount(1, 0))
			require.NoError(t, err)

			// assign $2
			err = interactor.MonthSetCategoryAmount(t, c.ctx, month.ID, category.ID, beans.NewAmount(2, 0))
			require.NoError(t, err)

			// check month and see if month category is assigned
			res, err := interactor.MonthGetOrCreate(t, c.ctx, testutils.NewMonthDate(t, "2022-04-01"))
			require.NoError(t, err)

			assert.Equal(t, 2, len(res.Categories)) // income + expense category
			findMonthCategory(t, res.Categories, category.ID, func(it beans.MonthCategoryWithDetails) {
				assert.Equal(t, beans.NewAmount(2, 0), it.Amount)
			})
		})
	})
}
