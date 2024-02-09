package specification

import "testing"

func DoTests(t *testing.T, interactor Interactor) {

	t.Run("account", func(t *testing.T) { testAccounts(t, interactor) })
	t.Run("budget", func(t *testing.T) { testBudgets(t, interactor) })
}
