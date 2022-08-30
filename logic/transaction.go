package logic

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type TransactionService struct {
	transactionRepository beans.TransactionRepository
}

func NewTransactionService(transactionRepository beans.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepository}
}

func (s *TransactionService) Create(ctx context.Context, data beans.TransactionCreate) (*beans.Transaction, error) {
	if err := data.ValidateAll(); err != nil {
		return nil, err
	}

	transactionID := beans.NewBeansID()
	err := s.transactionRepository.Create(ctx, transactionID, data.AccountID, data.Amount, data.Date, data.Notes)
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
