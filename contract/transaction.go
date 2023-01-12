package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
)

type TransactionContract struct {
	transactionRepository   beans.TransactionRepository
	accountRepository       beans.AccountRepository
	categoryRepository      beans.CategoryRepository
	monthCategoryRepository beans.MonthCategoryRepository
	monthRepository         beans.MonthRepository
}

func NewTransactionContract(
	transactionRepository beans.TransactionRepository,
	accountRepository beans.AccountRepository,
	categoryRepository beans.CategoryRepository,
	monthCategoryRepository beans.MonthCategoryRepository,
	monthRepository beans.MonthRepository,
) *TransactionContract {
	return &TransactionContract{transactionRepository, accountRepository, categoryRepository, monthCategoryRepository, monthRepository}
}

func (c *TransactionContract) Create(ctx context.Context, auth *beans.BudgetAuthContext, data beans.TransactionCreateParams) (*beans.Transaction, error) {
	if err := data.ValidateAll(); err != nil {
		return nil, err
	}

	account, err := c.accountRepository.Get(ctx, data.AccountID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.NewError(beans.EINVALID, "Invalid Account ID")
		} else {
			return nil, err
		}
	}
	if account.BudgetID != auth.BudgetID() {
		return nil, beans.NewError(beans.EINVALID, "Invalid Account ID")
	}

	if !data.CategoryID.Empty() {
		if _, err = c.categoryRepository.GetSingleForBudget(ctx, data.CategoryID, auth.BudgetID()); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return nil, beans.NewError(beans.EINVALID, "Invalid Category ID")
			} else {
				return nil, err
			}
		}

		month, err := c.monthRepository.GetOrCreate(ctx, auth.BudgetID(), beans.NewMonthDate(data.Date))
		if err != nil {
			return nil, err
		}

		if _, err := c.monthCategoryRepository.GetOrCreate(ctx, month.ID, data.CategoryID); err != nil {
			return nil, err
		}
	}

	transaction := &beans.Transaction{
		ID:         beans.NewBeansID(),
		AccountID:  data.AccountID,
		CategoryID: data.CategoryID,
		Amount:     data.Amount,
		Date:       data.Date,
		Notes:      data.Notes,

		Account: account,
	}
	err = c.transactionRepository.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (c *TransactionContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]*beans.Transaction, error) {
	return c.transactionRepository.GetForBudget(ctx, auth.BudgetID())
}
