package specification

import "testing"

func DoTests(t *testing.T, interactor Interactor) {

	t.Run("accounts", func(t *testing.T) { testAccounts(t, interactor) })
}
