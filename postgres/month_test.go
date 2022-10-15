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
	pool, container := StartPool(t)
	defer StopPool(t, container)

	monthRepository := postgres.NewMonthRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)
	budgetID2 := makeBudget(t, pool, "budget2", userID)

	cleanup := func() {
		pool.Exec(context.Background(), "truncate months cascade;")
	}

	t.Run("can create and get", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewDate(time.Now()), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month))

		res, err := monthRepository.GetByDate(context.Background(), budgetID, month.Date.Time)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(month, res))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewDate(time.Now()), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month))
		assertPgError(t, pgerrcode.UniqueViolation, monthRepository.Create(context.Background(), month))
	})

	t.Run("cannot create same month in same budget", func(t *testing.T) {
		defer cleanup()
		date := beans.NewDate(time.Now())
		month1 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		month2 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month1))
		assertPgError(t, pgerrcode.UniqueViolation, monthRepository.Create(context.Background(), month2))
	})

	t.Run("get month respects budget", func(t *testing.T) {
		defer cleanup()
		date := beans.NewDate(time.Now())
		month1 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		month2 := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID2}
		require.Nil(t, monthRepository.Create(context.Background(), month1))
		require.Nil(t, monthRepository.Create(context.Background(), month2))

		res, err := monthRepository.GetByDate(context.Background(), budgetID, date.Time)
		require.Nil(t, err)
		assert.Equal(t, month1.ID, res.ID)
	})

	t.Run("cannot get fictitious month", func(t *testing.T) {
		defer cleanup()
		_, err := monthRepository.GetByDate(context.Background(), budgetID, time.Now())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("get month ignores timezone", func(t *testing.T) {
		defer cleanup()
		date := beans.NewDate(time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC))
		loc, err := time.LoadLocation("America/New_York")
		require.Nil(t, err)
		month := &beans.Month{ID: beans.NewBeansID(), Date: date, BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month))

		res, err := monthRepository.GetByDate(context.Background(), budgetID, time.Date(2022, 05, 26, 23, 50, 0, 0, loc))
		require.Nil(t, err)
		assert.Equal(t, month.ID, res.ID)

		_, err = monthRepository.GetByDate(context.Background(), budgetID, time.Date(2022, 05, 25, 23, 50, 0, 0, loc))
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})
}