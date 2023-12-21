package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type TransactionContract struct {
	transactionRepository   beans.TransactionRepository
	accountRepository       beans.AccountRepository
	categoryRepository      beans.CategoryRepository
	monthCategoryRepository beans.MonthCategoryRepository
	monthRepository         beans.MonthRepository
	payeeRepository         beans.PayeeRepository
}

func NewTransactionContract(
	transactionRepository beans.TransactionRepository,
	accountRepository beans.AccountRepository,
	categoryRepository beans.CategoryRepository,
	monthCategoryRepository beans.MonthCategoryRepository,
	monthRepository beans.MonthRepository,
	payeeRepository beans.PayeeRepository,
) *TransactionContract {
	return &TransactionContract{transactionRepository, accountRepository, categoryRepository, monthCategoryRepository, monthRepository, payeeRepository}
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

		month, err := c.monthRepository.GetOrCreate(ctx, nil, auth.BudgetID(), beans.NewMonthDate(data.Date))
		if err != nil {
			return nil, err
		}

		if _, err := c.monthCategoryRepository.GetOrCreate(ctx, nil, month.ID, data.CategoryID); err != nil {
			return nil, err
		}
	}

	if err = c.validatePayee(ctx, auth, data.PayeeID); err != nil {
		return nil, err
	}

	transaction := &beans.Transaction{
		ID:         beans.NewBeansID(),
		AccountID:  data.AccountID,
		CategoryID: data.CategoryID,
		PayeeID:    data.PayeeID,
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

func (c *TransactionContract) Update(ctx context.Context, auth *beans.BudgetAuthContext, data beans.TransactionUpdateParams) error {
	if err := data.ValidateAll(); err != nil {
		return err
	}

	transaction, err := c.transactionRepository.Get(ctx, data.ID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.NewError(beans.EINVALID, "Invalid Transaction ID")
		}

		return err
	}

	if transaction.Account.BudgetID != auth.BudgetID() {
		return beans.NewError(beans.EINVALID, "Invalid Transaction ID")
	}

	account, err := c.accountRepository.Get(ctx, data.AccountID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.NewError(beans.EINVALID, "Invalid Account ID")
		}

		return err
	}
	if account.BudgetID != auth.BudgetID() {
		return beans.NewError(beans.EINVALID, "Invalid Account ID")
	}

	if !data.CategoryID.Empty() {
		if _, err = c.categoryRepository.GetSingleForBudget(ctx, data.CategoryID, auth.BudgetID()); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.NewError(beans.EINVALID, "Invalid Category ID")
			} else {
				return err
			}
		}

		month, err := c.monthRepository.GetOrCreate(ctx, nil, auth.BudgetID(), beans.NewMonthDate(data.Date))
		if err != nil {
			return err
		}

		if _, err := c.monthCategoryRepository.GetOrCreate(ctx, nil, month.ID, data.CategoryID); err != nil {
			return err
		}
	}

	if err = c.validatePayee(ctx, auth, data.PayeeID); err != nil {
		return err
	}

	transaction.AccountID = data.AccountID
	transaction.CategoryID = data.CategoryID
	transaction.PayeeID = data.PayeeID
	transaction.Amount = data.Amount
	transaction.Date = data.Date
	transaction.Notes = data.Notes

	err = c.transactionRepository.Update(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (c *TransactionContract) Delete(ctx context.Context, auth *beans.BudgetAuthContext, transactionIDs []beans.ID) error {
	return c.transactionRepository.Delete(ctx, auth.BudgetID(), transactionIDs)
}

func (c *TransactionContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]*beans.Transaction, error) {
	return c.transactionRepository.GetForBudget(ctx, auth.BudgetID())
}

func (c *TransactionContract) validatePayee(ctx context.Context, auth *beans.BudgetAuthContext, payeeID beans.ID) error {
	if !payeeID.Empty() {
		payee, err := c.payeeRepository.Get(ctx, payeeID)

		if err != nil && errors.Is(err, beans.ErrorNotFound) {
			return beans.NewError(beans.EINVALID, "Invalid Payee ID")
		} else if err != nil {
			return err
		} else if err == nil && payee.BudgetID != auth.BudgetID() {
			return beans.NewError(beans.EINVALID, "Invalid Payee ID")
		}
	}

	return nil
}
