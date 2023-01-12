package testutils

import (
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/require"
)

func BudgetAuthContext(t testing.TB, userID beans.ID, budget *beans.Budget) *beans.BudgetAuthContext {
	auth, err := beans.NewBudgetAuthContext(beans.NewAuthContext(userID), budget)
	require.Nil(t, err)
	return auth
}
