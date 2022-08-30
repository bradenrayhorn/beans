package logic

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type accountService struct {
	accountRepository beans.AccountRepository
}

func NewAccountService(ar beans.AccountRepository) *accountService {
	return &accountService{accountRepository: ar}
}

func (s *accountService) Create(ctx context.Context, name beans.Name, budgetID beans.ID) (*beans.Account, error) {
	if err := beans.ValidateFields(beans.Field("Account name", name)); err != nil {
		return nil, err
	}

	accountID := beans.NewBeansID()
	if err := s.accountRepository.Create(ctx, accountID, name, budgetID); err != nil {
		return nil, err
	}

	return &beans.Account{
		ID:       accountID,
		Name:     name,
		BudgetID: budgetID,
	}, nil
}
