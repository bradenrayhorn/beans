package logic_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/require"
)

func TestCreateBudget(t *testing.T) {
	budgetRepository := mocks.NewMockBudgetRepository()
	svc := logic.NewBudgetService(budgetRepository)

	t.Run("name is required", func(t *testing.T) {
		_, err := svc.CreateBudget(context.Background(), beans.Name(""), beans.UserID(beans.NewBeansID()))
		testutils.AssertError(t, err, "Budget name is required.")
	})

	t.Run("can create budget", func(t *testing.T) {
		budget, err := svc.CreateBudget(context.Background(), beans.Name("budget1"), beans.UserID(beans.NewBeansID()))
		require.Nil(t, err)

		require.Equal(t, beans.Name("budget1"), budget.Name)
		require.False(t, budget.ID.Empty())
	})
}
