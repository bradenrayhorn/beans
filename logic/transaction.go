package logic

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
)

type TransactionService struct {
	transactionRepository   beans.TransactionRepository
	accountRepository       beans.AccountRepository
	categoryRepository      beans.CategoryRepository
	monthCategoryRepository beans.MonthCategoryRepository
	monthRepository         beans.MonthRepository
}

func NewTransactionService(
	transactionRepository beans.TransactionRepository,
	accountRepository beans.AccountRepository,
	categoryRepository beans.CategoryRepository,
	monthCategoryRepository beans.MonthCategoryRepository,
	monthRepository beans.MonthRepository,
) *TransactionService {
	return &TransactionService{transactionRepository, accountRepository, categoryRepository, monthCategoryRepository, monthRepository}
}

func (s *TransactionService) Create(ctx context.Context, activeBudget *beans.Budget, data beans.TransactionCreate) (*beans.Transaction, error) {
	if err := data.ValidateAll(); err != nil {
		return nil, err
	}

	account, err := s.accountRepository.Get(ctx, data.AccountID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.NewError(beans.EINVALID, "Invalid Account ID")
		} else {
			return nil, err
		}
	}
	if account.BudgetID != activeBudget.ID {
		return nil, beans.NewError(beans.EINVALID, "Invalid Account ID")
	}

	if !data.CategoryID.Empty() {
		if _, err = s.categoryRepository.GetSingleForBudget(ctx, data.CategoryID, activeBudget.ID); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return nil, beans.NewError(beans.EINVALID, "Invalid Category ID")
			} else {
				return nil, err
			}
		}

		month, err := s.monthRepository.GetOrCreate(ctx, activeBudget.ID, beans.NewMonthDate(data.Date))
		if err != nil {
			return nil, err
		}

		if _, err := s.monthCategoryRepository.GetOrCreate(ctx, month.ID, data.CategoryID); err != nil {
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
	err = s.transactionRepository.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
