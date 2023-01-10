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
		pool.Exec(context.Background(), "truncate months cascade;")
	}

	t.Run("can create and get", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month))

		res, err := monthRepository.Get(context.Background(), month.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("create respects tx", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewMonthDate(beans.NewDate(time.Now().AddDate(0, 1, 0))), BudgetID: budgetID}

		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer tx.Rollback(context.Background())

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

	t.Run("get or create respects budget", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		monthBudget2 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID2}
		require.Nil(t, monthRepository.Create(context.Background(), nil, monthBudget2))

		month, err := monthRepository.GetOrCreate(context.Background(), budgetID, date)
		require.Nil(t, err)

		assert.NotEqual(t, month.ID, monthBudget2.ID)
	})

	t.Run("get or create returns existing month", func(t *testing.T) {
		defer cleanup()
		date := beans.NewMonthDate(beans.NewDate(time.Now()))

		existingMonth := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, existingMonth))

		month, err := monthRepository.GetOrCreate(context.Background(), budgetID, date)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(month, existingMonth))
	})

	t.Run("get or create creates new month", func(t *testing.T) {
		defer cleanup()
		date1 := testutils.NewMonthDate(t, "2022-05-01")
		date2 := testutils.NewMonthDate(t, "2022-06-01")

		existingMonth := &beans.Month{ID: beans.NewBeansID(), Date: date1, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, existingMonth))

		month, err := monthRepository.GetOrCreate(context.Background(), budgetID, date2)
		require.Nil(t, err)

		assert.NotEqual(t, existingMonth.ID, month.ID)
		assert.True(t, reflect.DeepEqual(
			&beans.Month{
				ID:       month.ID,
				Date:     date2,
				BudgetID: budgetID,
			},
			month,
		))
	})

	t.Run("cannot get fictitious month by id", func(t *testing.T) {
		defer cleanup()
		_, err := monthRepository.Get(context.Background(), beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get latest month", func(t *testing.T) {
		defer cleanup()
		month1 := &beans.Month{ID: beans.NewBeansID(), Date: testutils.NewMonthDate(t, "2022-05-01"), BudgetID: budgetID}
		month2 := &beans.Month{ID: beans.NewBeansID(), Date: testutils.NewMonthDate(t, "2022-07-01"), BudgetID: budgetID}
		month3 := &beans.Month{ID: beans.NewBeansID(), Date: testutils.NewMonthDate(t, "2022-03-01"), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), nil, month1))
		require.Nil(t, monthRepository.Create(context.Background(), nil, month2))
		require.Nil(t, monthRepository.Create(context.Background(), nil, month3))

		res, err := monthRepository.GetLatest(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Equal(t, month2.ID, res.ID)
	})

	t.Run("can get latest month when none exists", func(t *testing.T) {
		defer cleanup()
		_, err := monthRepository.GetLatest(context.Background(), budgetID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})
}
