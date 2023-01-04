package logic_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBudget(t *testing.T) {
	budgetRepository := mocks.NewMockBudgetRepository()
	monthRepository := mocks.NewMockMonthRepository()
	txManager := mocks.NewMockTxManager()
	svc := logic.NewBudgetService(txManager, budgetRepository, monthRepository)

	t.Run("name is required", func(t *testing.T) {
		_, err := svc.CreateBudget(context.Background(), beans.Name(""), beans.UserID(beans.NewBeansID()))
		testutils.AssertError(t, err, "Budget name is required.")
	})

	t.Run("can create budget", func(t *testing.T) {
		tx := mocks.NewMockTx()
		txManager.CreateFunc.PushReturn(tx, nil)

		budget, err := svc.CreateBudget(context.Background(), beans.Name("budget1"), beans.UserID(beans.NewBeansID()))
		require.Nil(t, err)

		require.Equal(t, beans.Name("budget1"), budget.Name)
		require.False(t, budget.ID.Empty())

		createdMonth := monthRepository.CreateFunc.History()[0].Arg1
		assert.Equal(t, budget.ID, createdMonth.BudgetID)
		assert.Equal(t, beans.NewDate(beans.NormalizeMonth(time.Now())), createdMonth.Date)
	})
}
