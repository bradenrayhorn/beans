package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonth(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	txManager := postgres.NewTxManager(pool)
	monthRepository := postgres.NewMonthRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID
	budgetID2 := testutils.MakeBudget(t, pool, "budget2", userID).ID

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate months cascade;")
	}

	t.Run("can create and get", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{
			ID:        beans.NewBeansID(),
			Date:      beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))),
			BudgetID:  budgetID,
			Carryover: beans.NewAmount(5, 0),
		}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))

		res, err := monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("can create with no carryover", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{
			ID:       beans.NewBeansID(),
			Date:     beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))),
			BudgetID: budgetID,
		}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))

		res, err := monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)

		// carryover should have been initialized to 0
		month.Carryover = beans.NewAmount(0, 0)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("create respects tx", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))), BudgetID: budgetID}

		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		require.Nil(t, monthRepository.Create(context.Background(), tx, month))

		_, err = monthRepository.Get(context.Background(), month.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		require.Nil(t, tx.Commit(context.Background()))

		_, err = monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewMonthDate(beans.NewDate(time.Now())), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))
		assertPgError(t, pgerrcode.UniqueViolation, monthRepository.Create(context.Background(), nil, month))
	})

	t.Run("cannot create same month in same budget", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))
		month1 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		month2 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month1))
		assertPgError(t, pgerrcode.UniqueViolation, monthRepository.Create(context.Background(), nil, month2))
	})

	t.Run("can update month carryover", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		month := &beans.Month{
			ID:        beans.NewBeansID(),
			Date:      date,
			BudgetID:  budgetID,
			Carryover: beans.NewAmount(0, 0),
		}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))

		month.Carryover = beans.NewAmount(5, 0)
		require.Nil(t, monthRepository.Update(context.Background(), month))

		res, err := monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("can update month carryover to nil", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		month := &beans.Month{
			ID:        beans.NewBeansID(),
			Date:      date,
			BudgetID:  budgetID,
			Carryover: beans.NewAmount(5, 0),
		}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))

		month.Carryover = beans.NewEmptyAmount()
		require.Nil(t, monthRepository.Update(context.Background(), month))

		res, err := monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)

		// carryover should have been reset to 0
		month.Carryover = beans.NewAmount(0, 0)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("get or create respects budget", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		monthBudget2 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID2}
		require.Nil(t, monthRepository.Create(context.Background(), nil, monthBudget2))

		month, err := monthRepository.GetOrCreate(context.Background(), nil, budgetID, date)
		require.Nil(t, err)

		assert.NotEqual(t, month.ID, monthBudget2.ID)
	})

	t.Run("get or create returns existing month", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		existingMonth := &beans.Month{
			ID:        beans.NewBeansID(),
			Date:      date,
			BudgetID:  budgetID,
			Carryover: beans.NewAmount(0, 0),
		}
		require.Nil(t, monthRepository.Create(context.Background(), nil, existingMonth))

		month, err := monthRepository.GetOrCreate(context.Background(), nil, budgetID, date)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(month, existingMonth))
	})

	t.Run("get or create creates new month", func(t *testing.T) {
		defer cleanup()
		date1 := testutils.NewMonthDate(t, "2022-05-01")
		date2 := testutils.NewMonthDate(t, "2022-06-01")

		existingMonth := &beans.Month{ID: beans.NewBeansID(), Date: date1, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, existingMonth))

		month, err := monthRepository.GetOrCreate(context.Background(), nil, budgetID, date2)
		require.Nil(t, err)

		assert.NotEqual(t, existingMonth.ID, month.ID)
		assert.True(t, reflect.DeepEqual(
			&beans.Month{
				ID:        month.ID,
				Date:      date2,
				BudgetID:  budgetID,
				Carryover: beans.NewAmount(0, 0),
			},
			month,
		))
	})

	t.Run("get or create respects tx", func(t *testing.T) {
		defer cleanup()
		date := testutils.NewMonthDate(t, "2022-05-01")

		// make transaction
		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		// get or create but do not commit
		month1, err := monthRepository.GetOrCreate(context.Background(), tx, budgetID, date)
		require.Nil(t, err)

		// try to find month, should fail
		_, err = monthRepository.Get(context.Background(), month1.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		// commit
		require.Nil(t, tx.Commit(context.Background()))

		// try to find month, should succeed
		_, err = monthRepository.Get(context.Background(), month1.ID)
		require.Nil(t, err)
	})

	t.Run("cannot get fictitious month by id", func(t *testing.T) {
		defer cleanup()
		_, err := monthRepository.Get(context.Background(), beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get months in budget", func(t *testing.T) {
		defer cleanup()
		month1 := &beans.Month{ID: beans.NewBeansID(), Date: testutils.NewMonthDate(t, "2022-05-01"), BudgetID: budgetID}
		month2 := &beans.Month{ID: beans.NewBeansID(), Date: testutils.NewMonthDate(t, "2022-07-01"), BudgetID: budgetID2}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month1))
		require.Nil(t, monthRepository.Create(context.Background(), nil, month2))

		res, err := monthRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, month1.ID, res[0].ID)
	})
}
