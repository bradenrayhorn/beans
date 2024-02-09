package datasource

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonthRepository(t *testing.T, ds beans.DataSource) {
	factory := testutils.Factory(t, ds)

	monthRepository := ds.MonthRepository()
	ctx := context.Background()

	t.Run("can create and get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month := beans.Month{
			ID:        beans.NewBeansID(),
			Date:      beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))),
			BudgetID:  budget.ID,
			Carryover: beans.NewAmount(5, 0),
		}
		require.Nil(t, monthRepository.Create(ctx, nil, month))

		res, err := monthRepository.Get(ctx, budget.ID, month.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("can create with no carryover", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month := beans.Month{
			ID:       beans.NewBeansID(),
			Date:     beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))),
			BudgetID: budget.ID,
		}
		require.Nil(t, monthRepository.Create(ctx, nil, month))

		res, err := monthRepository.Get(ctx, budget.ID, month.ID)
		require.Nil(t, err)

		// carryover should have been initialized to 0
		month.Carryover = beans.NewAmount(0, 0)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("create respects tx", func(t *testing.T) {
		txManager := ds.TxManager()
		budget, _ := factory.MakeBudgetAndUser()

		month := beans.Month{ID: beans.NewBeansID(), Date: beans.NewMonthDate(beans.NewDate(time.Now())), BudgetID: budget.ID}

		tx, err := txManager.Create(ctx)
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		require.Nil(t, monthRepository.Create(ctx, tx, month))

		_, err = monthRepository.Get(ctx, budget.ID, month.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		require.Nil(t, tx.Commit(ctx))

		_, err = monthRepository.Get(ctx, budget.ID, month.ID)
		require.Nil(t, err)
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month := factory.Month(beans.Month{BudgetID: budget.ID})
		assert.NotNil(t, monthRepository.Create(ctx, nil, month))
	})

	t.Run("cannot create same month in same budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		date := beans.NewMonthDate(beans.NewDate(time.Now()))
		month1 := beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budget.ID}
		month2 := beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budget.ID}

		assert.Nil(t, monthRepository.Create(ctx, nil, month1))
		assert.NotNil(t, monthRepository.Create(ctx, nil, month2))
	})

	t.Run("can update month carryover", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		month := factory.Month(beans.Month{BudgetID: budget.ID, Carryover: beans.NewAmount(0, 0)})

		month.Carryover = beans.NewAmount(5, 0)
		require.Nil(t, monthRepository.Update(ctx, month))

		res, err := monthRepository.Get(ctx, budget.ID, month.ID)
		require.Nil(t, err)
		assert.Equal(t, beans.NewAmount(5, 0), res.Carryover)
	})

	t.Run("can update month carryover to nil", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID, Carryover: beans.NewAmount(5, 0)})

		// update month to have an empty carryover
		month.Carryover = beans.NewEmptyAmount()
		require.Nil(t, monthRepository.Update(ctx, month))

		// get month, carryover should have been reset to 0
		res, err := monthRepository.Get(ctx, budget.ID, month.ID)
		require.Nil(t, err)

		assert.Equal(t, beans.NewAmount(0, 0), res.Carryover)
	})

	t.Run("get or create respects budget", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		budget2Month := factory.Month(beans.Month{BudgetID: budget2.ID, Date: date})

		month, err := monthRepository.GetOrCreate(ctx, nil, budget1.ID, date)
		require.Nil(t, err)

		assert.NotEqual(t, month.ID, budget2Month.ID)
	})

	t.Run("get or create returns existing month", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		existingMonth := factory.Month(beans.Month{BudgetID: budget.ID, Date: date})

		month, err := monthRepository.GetOrCreate(ctx, nil, budget.ID, date)
		require.Nil(t, err)

		assert.Equal(t, existingMonth.ID, month.ID)
	})

	t.Run("get or create creates new month", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		date1 := testutils.NewMonthDate(t, "2022-05-01")
		date2 := testutils.NewMonthDate(t, "2022-06-01")

		existingMonth := factory.Month(beans.Month{BudgetID: budget.ID, Date: date1})

		month, err := monthRepository.GetOrCreate(ctx, nil, budget.ID, date2)
		require.Nil(t, err)

		assert.NotEqual(t, existingMonth.ID, month.ID)
	})

	t.Run("get or create respects tx", func(t *testing.T) {
		txManager := ds.TxManager()
		budget, _ := factory.MakeBudgetAndUser()

		date := testutils.NewMonthDate(t, "2022-05-01")

		// make transaction
		tx, err := txManager.Create(ctx)
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// get or create but do not commit
		month1, err := monthRepository.GetOrCreate(ctx, tx, budget.ID, date)
		require.Nil(t, err)

		// try to find month, should fail
		_, err = monthRepository.Get(ctx, budget.ID, month1.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		// commit
		require.Nil(t, tx.Commit(ctx))

		// try to find month, should succeed
		_, err = monthRepository.Get(ctx, budget.ID, month1.ID)
		require.Nil(t, err)
	})

	t.Run("cannot get fictitious month by id", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		_, err := monthRepository.Get(ctx, budget.ID, beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("cannot get month from other budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		month := factory.Month(beans.Month{BudgetID: budget.ID})

		_, err := monthRepository.Get(ctx, budget2.ID, month.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get months in budget", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		expected := factory.Month(beans.Month{BudgetID: budget1.ID})
		factory.Month(beans.Month{BudgetID: budget2.ID})

		res, err := monthRepository.GetForBudget(ctx, budget1.ID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, expected.ID, res[0].ID)
	})
}
