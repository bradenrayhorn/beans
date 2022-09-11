package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/require"
)

func TestMonth(t *testing.T) {
	pool, container := StartPool(t)
	defer StopPool(t, container)

	monthRepository := postgres.NewMonthRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)

	cleanup := func() {
		pool.Exec(context.Background(), "truncate months;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewDate(time.Now()), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		month := &beans.Month{ID: beans.NewBeansID(), Date: beans.NewDate(time.Now()), BudgetID: budgetID}
		require.Nil(t, monthRepository.Create(context.Background(), month))
		assertPgError(t, pgerrcode.UniqueViolation, monthRepository.Create(context.Background(), month))
	})
}
