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
	categoryRepository := mocks.NewMockCategoryRepository()
	txManager := mocks.NewMockTxManager()
	svc := logic.NewBudgetService(txManager, budgetRepository, monthRepository, categoryRepository)

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

		createdMonth := monthRepository.CreateFunc.History()[0].Arg2
		assert.Equal(t, budget.ID, createdMonth.BudgetID)
		assert.Equal(t, beans.NewDate(beans.NormalizeMonth(time.Now())), createdMonth.Date)

		// category was created
		group := categoryRepository.CreateGroupFunc.History()[0].Arg2
		assert.Equal(t, budget.ID, group.BudgetID)
		assert.Equal(t, "Income", string(group.Name))

		category := categoryRepository.CreateFunc.History()[0].Arg2
		assert.Equal(t, budget.ID, category.BudgetID)
		assert.Equal(t, group.ID, category.GroupID)
		assert.Equal(t, "Income", string(category.Name))
		assert.Equal(t, true, category.IsIncome)

		// transaction was committed
		assert.Len(t, tx.CommitFunc.History(), 1)
	})
}
