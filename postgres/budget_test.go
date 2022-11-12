package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgets(t *testing.T) {
	pool, stop := testutils.StartPool(t)
	defer stop()

	budgetRepository := postgres.NewBudgetRepository(pool)
	monthRepository := postgres.NewMonthRepository(pool)

	userID := makeUser(t, pool, "user")

	t.Run("can create and get budget", func(t *testing.T) {
		defer pool.Exec(context.Background(), "truncate budgets;")
		budgetID := beans.NewBeansID()
		err := budgetRepository.Create(context.Background(), budgetID, "Budget1", userID, time.Now())
		require.Nil(t, err)

		month, err := monthRepository.GetByDate(context.Background(), budgetID, beans.NormalizeMonth(time.Now()))
		require.Nil(t, err)
		assert.Equal(t, budgetID, month.BudgetID)
		assert.Equal(t, beans.NewDate(beans.NormalizeMonth(time.Now())), month.Date)

		budget, err := budgetRepository.Get(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Equal(t, budgetID, budget.ID)
		assert.Equal(t, "Budget1", string(budget.Name))
		assert.Equal(t, []beans.UserID{userID}, budget.UserIDs)
	})
}
