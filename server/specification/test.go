package specification

import "testing"

func DoTests(t *testing.T, interactor Interactor) {

	t.Run("account", func(t *testing.T) {
		t.Parallel()
		testAccount(t, interactor)
	})
	t.Run("budget", func(t *testing.T) {
		t.Parallel()
		testBudget(t, interactor)
	})
	t.Run("category", func(t *testing.T) {
		t.Parallel()
		testCategory(t, interactor)
	})
	t.Run("month", func(t *testing.T) {
		t.Parallel()
		testMonth(t, interactor)
	})
	t.Run("payee", func(t *testing.T) {
		t.Parallel()
		testPayee(t, interactor)
	})
	t.Run("transaction", func(t *testing.T) {
		t.Parallel()
		testTransaction(t, interactor)
	})
	t.Run("user", func(t *testing.T) {
		t.Parallel()
		testUser(t, interactor)
	})
}
