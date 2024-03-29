package testutils

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/require"
)

func BudgetAuthContext(t testing.TB, userID beans.ID, budget beans.Budget) *beans.BudgetAuthContext {
	auth, err := beans.NewBudgetAuthContext(beans.NewAuthContext(userID, beans.SessionID("1234")), budget)
	require.Nil(t, err)
	return auth
}
