package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBudgets(t *testing.T, interactor Interactor) {

	t.Run("create", func(t *testing.T) {
		t.Run("cannot create with invalid name", func(t *testing.T) {
			c := makeUser(t, interactor)

			_, err := interactor.BudgetCreate(t, c.ctx, "")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUser(t, interactor)

			// create budget
			budgetID, err := interactor.BudgetCreate(t, c.ctx, "New Budget")
			require.NoError(t, err)

			// get budget
			budget, err := interactor.BudgetGet(t, c.ctx, budgetID)
			require.NoError(t, err)

			assert.Equal(t, beans.Name("New Budget"), budget.Name)

			// get budget categories
			groups, err := interactor.CategoryGetAll(t, Context{SessionID: c.sessionID, BudgetID: budget.ID})
			require.NoError(t, err)

			// check income group is created
			require.Len(t, groups, 1)
			incomeGroup := groups[0]

			assert.False(t, incomeGroup.ID.Empty())
			assert.Equal(t, beans.Name("Income"), incomeGroup.Name)
			assert.Equal(t, true, incomeGroup.IsIncome)

			// check income category is created
			require.Len(t, incomeGroup.Categories, 1)
			incomeCategory := incomeGroup.Categories[0]

			assert.False(t, incomeCategory.ID.Empty())
			assert.Equal(t, beans.Name("Income"), incomeCategory.Name)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get non-existent budget", func(t *testing.T) {
			c := makeUser(t, interactor)

			_, err := interactor.BudgetGet(t, c.ctx, beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot budget of another user", func(t *testing.T) {
			c1 := makeUser(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			_, err := interactor.BudgetGet(t, c1.ctx, c2.budget.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all", func(t *testing.T) {
			c := makeUser(t, interactor)

			// this budget should show in the response
			id, err := interactor.BudgetCreate(t, c.ctx, "New Budget")
			require.NoError(t, err)

			// this budget should not show in the response
			makeUserAndBudget(t, interactor)

			// get budgets the user has access to
			budgets, err := interactor.BudgetGetAll(t, c.ctx)
			require.NoError(t, err)

			require.Len(t, budgets, 1)
			budget := budgets[0]
			require.Equal(t, id, budget.ID)
			require.Equal(t, beans.Name("New Budget"), budget.Name)
		})
	})
}
