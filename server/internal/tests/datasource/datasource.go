package datasource

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
)

func DoTestDatasource(t *testing.T, ds beans.DataSource) {

	t.Run("account", func(t *testing.T) { testAccount(t, ds) })
	t.Run("budget", func(t *testing.T) { testBudget(t, ds) })
	t.Run("category", func(t *testing.T) { testCategory(t, ds) })
	t.Run("month", func(t *testing.T) { testMonth(t, ds) })
	t.Run("month category", func(t *testing.T) { testMonthCategory(t, ds) })
	t.Run("payee", func(t *testing.T) { testPayee(t, ds) })
	t.Run("transaction", func(t *testing.T) { testTransaction(t, ds) })
	t.Run("user", func(t *testing.T) { testUser(t, ds) })
}
