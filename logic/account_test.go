package logic_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	accountRepository := mocks.NewMockAccountRepository()
	svc := logic.NewAccountService(accountRepository)

	t.Run("name is required", func(t *testing.T) {
		_, err := svc.Create(context.Background(), beans.Name(""), beans.NewBeansID())
		testutils.AssertError(t, err, "Account name is required.")
	})

	t.Run("can create account", func(t *testing.T) {
		budgetID := beans.NewBeansID()
		account, err := svc.Create(context.Background(), beans.Name("my account"), budgetID)
		require.Nil(t, err)

		require.Equal(t, beans.Name("my account"), account.Name)
		require.Equal(t, budgetID, account.BudgetID)
		require.False(t, account.ID.Empty())
	})
}
