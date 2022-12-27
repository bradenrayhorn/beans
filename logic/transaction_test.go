package logic_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	budget := &beans.Budget{
		ID:   beans.NewBeansID(),
		Name: "Budget1",
	}
	month := &beans.Month{
		ID:       beans.NewBeansID(),
		BudgetID: budget.ID,
		Date:     beans.NewDate(time.Now()),
	}
	account := &beans.Account{
		ID:       beans.NewBeansID(),
		Name:     "Account1",
		BudgetID: budget.ID,
	}
	categoryGroup := &beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		Name:     "Group1",
		BudgetID: budget.ID,
	}
	category := &beans.Category{
		ID:       beans.NewBeansID(),
		Name:     "Category1",
		BudgetID: budget.ID,
		GroupID:  categoryGroup.ID,
	}

	transactionRepository := mocks.NewMockTransactionRepository()
	accountRepository := mocks.NewMockAccountRepository()
	categoryRepository := mocks.NewMockCategoryRepository()
	monthService := mocks.NewMockMonthService()
	monthCategoryService := mocks.NewMockMonthCategoryService()

	monthService.GetOrCreateFunc.SetDefaultReturn(nil, errors.New("invalid"))
	monthCategoryService.CreateIfNotExistsFunc.SetDefaultReturn(errors.New("invalid"))
	categoryRepository.GetSingleForBudgetFunc.SetDefaultReturn(nil, errors.New("invalid"))
	svc := logic.NewTransactionService(transactionRepository, accountRepository, categoryRepository, monthService, monthCategoryService)

	t.Run("fields are required", func(t *testing.T) {
		_, err := svc.Create(context.Background(), budget, beans.TransactionCreate{})
		testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
	})

	t.Run("cannot create transaction with amount more than 2 decimals", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, -3),
			Date:      beans.NewDate(time.Now()),
		}
		_, err := svc.Create(context.Background(), budget, c)
		testutils.AssertError(t, err, "Amount must have at most 2 decimal points.")
	})

	t.Run("can create full", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(10, 1),
			Date:       testutils.NewDate(t, "2022-06-07"),
			Notes:      beans.NewTransactionNotes("My Notes"),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)
		categoryRepository.GetSingleForBudgetFunc.PushReturn(category, nil)
		monthService.GetOrCreateFunc.PushReturn(month, nil)
		monthCategoryService.CreateIfNotExistsFunc.PushReturn(nil)

		transaction, err := svc.Create(context.Background(), budget, c)
		require.Nil(t, err)
		require.Equal(t, c.AccountID, transaction.AccountID)
		require.Equal(t, c.CategoryID, transaction.CategoryID)
		require.Equal(t, c.Amount, transaction.Amount)
		require.Equal(t, c.Date, transaction.Date)
		require.Equal(t, c.Notes, transaction.Notes)
		assert.True(t, reflect.DeepEqual(account, transaction.Account))

		assert.Equal(t, testutils.NewDate(t, "2022-06-01").Time, monthService.GetOrCreateFunc.History()[0].Arg2)
	})

	t.Run("can create minimum", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)

		transaction, err := svc.Create(context.Background(), budget, c)
		require.Nil(t, err)
		require.Equal(t, c.AccountID, transaction.AccountID)
		require.Equal(t, c.Amount, transaction.Amount)
		require.Equal(t, c.Date, transaction.Date)
		assert.True(t, reflect.DeepEqual(account, transaction.Account))
	})

	t.Run("cannot create after account error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My Notes"),
		}

		accountRepository.GetFunc.SetDefaultReturn(nil, errors.New("account not found"))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "account not found")
	})

	t.Run("translates account check not found error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My Notes"),
		}

		accountRepository.GetFunc.SetDefaultReturn(nil, beans.WrapError(errors.New("not found"), beans.ErrorNotFound))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "Invalid Account ID")
	})

	t.Run("cannot create with account from other budget", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID: beans.NewBeansID(),
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My notes"),
		}
		badAccount := &beans.Account{
			ID:       c.AccountID,
			Name:     "bad account",
			BudgetID: beans.NewBeansID(),
		}
		accountRepository.GetFunc.SetDefaultReturn(badAccount, nil)

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		testutils.AssertError(t, err, "Invalid Account ID")
	})

	t.Run("cannot create after category error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(10, 1),
			Date:       beans.NewDate(time.Now()),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)
		categoryRepository.GetSingleForBudgetFunc.PushReturn(nil, errors.New("some error"))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "some error")
	})

	t.Run("translates category check not found error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(10, 1),
			Date:       beans.NewDate(time.Now()),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)
		categoryRepository.GetSingleForBudgetFunc.SetDefaultReturn(nil, beans.WrapError(errors.New("not found"), beans.ErrorNotFound))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "Invalid Category ID")
	})

	t.Run("cannot create after get month error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(10, 1),
			Date:       beans.NewDate(time.Now()),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)
		categoryRepository.GetSingleForBudgetFunc.PushReturn(category, nil)
		monthService.GetOrCreateFunc.PushReturn(nil, errors.New("some error"))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "some error")
	})

	t.Run("cannot create after create month category error", func(t *testing.T) {
		c := beans.TransactionCreate{
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(10, 1),
			Date:       beans.NewDate(time.Now()),
		}

		accountRepository.GetFunc.SetDefaultReturn(account, nil)
		categoryRepository.GetSingleForBudgetFunc.PushReturn(category, nil)
		monthService.GetOrCreateFunc.PushReturn(month, nil)
		monthCategoryService.CreateIfNotExistsFunc.PushReturn(errors.New("some error"))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "some error")
	})
}
