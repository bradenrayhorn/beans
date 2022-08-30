package logic

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type TransactionService struct {
	transactionRepository beans.TransactionRepository
	accountRepository     beans.AccountRepository
}

func NewTransactionService(transactionRepository beans.TransactionRepository, accountRepository beans.AccountRepository) *TransactionService {
	return &TransactionService{transactionRepository, accountRepository}
}

func (s *TransactionService) Create(ctx context.Context, activeBudget *beans.Budget, data beans.TransactionCreate) (*beans.Transaction, error) {
	if err := data.ValidateAll(); err != nil {
		return nil, err
	}

	account, err := s.accountRepository.Get(ctx, data.AccountID)
	if err != nil {
		return nil, err
	}
	if account.BudgetID != activeBudget.ID {
		return nil, beans.NewError(beans.EINVALID, "Invalid Account ID")
	}

	transactionID := beans.NewBeansID()
	err = s.transactionRepository.Create(ctx, transactionID, data.AccountID, data.Amount, data.Date, data.Notes)
	if err != nil {
		return nil, err
	}

	return &beans.Transaction{
		ID:        transactionID,
		AccountID: data.AccountID,
		Amount:    data.Amount,
		Date:      data.Date,
		Notes:     data.Notes,
	}, nil
}
