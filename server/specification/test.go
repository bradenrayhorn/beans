package specification

import "testing"

func DoTests(t *testing.T, interactor Interactor) {

	t.Run("account", func(t *testing.T) { testAccount(t, interactor) })
	t.Run("budget", func(t *testing.T) { testBudget(t, interactor) })
	t.Run("category", func(t *testing.T) { testCategory(t, interactor) })
	t.Run("month", func(t *testing.T) { testMonth(t, interactor) })
	t.Run("payee", func(t *testing.T) { testPayee(t, interactor) })
	t.Run("transaction", func(t *testing.T) { testTransaction(t, interactor) })
	t.Run("user", func(t *testing.T) { testUser(t, interactor) })
}
