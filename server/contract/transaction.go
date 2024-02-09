package contract

import (
	"context"
	"errors"
	"fmt"

	"github.com/bradenrayhorn/beans/server/beans"
)

type transactionContract struct{ contract }

var _ beans.TransactionContract = (*transactionContract)(nil)

func (c *transactionContract) Create(ctx context.Context, auth *beans.BudgetAuthContext, data beans.TransactionCreateParams) (beans.ID, error) {
	if err := data.ValidateAll(); err != nil {
		return beans.EmptyID(), err
	}

	account, err := c.ds().AccountRepository().Get(ctx, data.AccountID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.EmptyID(), beans.NewError(beans.EINVALID, "Invalid Account ID")
		} else {
			return beans.EmptyID(), err
		}
	}
	if account.BudgetID != auth.BudgetID() {
		return beans.EmptyID(), beans.NewError(beans.EINVALID, "Invalid Account ID")
	}

	if !data.CategoryID.Empty() {
		if _, err = c.ds().CategoryRepository().GetSingleForBudget(ctx, data.CategoryID, auth.BudgetID()); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.EmptyID(), beans.NewError(beans.EINVALID, "Invalid Category ID")
			} else {
				return beans.EmptyID(), err
			}
		}

		month, err := c.ds().MonthRepository().GetOrCreate(ctx, nil, auth.BudgetID(), beans.NewMonthDate(data.Date))
		if err != nil {
			return beans.EmptyID(), err
		}

		if _, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, nil, month.ID, data.CategoryID); err != nil {
			return beans.EmptyID(), err
		}
	}

	if err = c.validatePayee(ctx, auth, data.PayeeID); err != nil {
		return beans.EmptyID(), err
	}

	transaction := beans.Transaction{
		ID:         beans.NewBeansID(),
		AccountID:  data.AccountID,
		CategoryID: data.CategoryID,
		PayeeID:    data.PayeeID,
		Amount:     data.Amount,
		Date:       data.Date,
		Notes:      data.Notes,
	}
	err = c.ds().TransactionRepository().Create(ctx, transaction)
	if err != nil {
		return beans.EmptyID(), err
	}

	return transaction.ID, nil
}

func (c *transactionContract) Update(ctx context.Context, auth *beans.BudgetAuthContext, data beans.TransactionUpdateParams) error {
	if err := data.ValidateAll(); err != nil {
		return err
	}

	transaction, err := c.ds().TransactionRepository().Get(ctx, data.ID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.NewError(beans.EINVALID, "Invalid Transaction ID")
		}

		return err
	}

	currentAccount, err := c.ds().AccountRepository().Get(ctx, transaction.AccountID)
	if err != nil {
		return fmt.Errorf("could not find related account: %w", err)
	}
	if currentAccount.BudgetID != auth.BudgetID() {
		return beans.NewError(beans.EINVALID, "Invalid Transaction ID")
	}

	account, err := c.ds().AccountRepository().Get(ctx, data.AccountID)
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
		if _, err = c.ds().CategoryRepository().GetSingleForBudget(ctx, data.CategoryID, auth.BudgetID()); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.NewError(beans.EINVALID, "Invalid Category ID")
			} else {
				return err
			}
		}

		month, err := c.ds().MonthRepository().GetOrCreate(ctx, nil, auth.BudgetID(), beans.NewMonthDate(data.Date))
		if err != nil {
			return err
		}

		if _, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, nil, month.ID, data.CategoryID); err != nil {
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

	err = c.ds().TransactionRepository().Update(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (c *transactionContract) Delete(ctx context.Context, auth *beans.BudgetAuthContext, transactionIDs []beans.ID) error {
	return c.ds().TransactionRepository().Delete(ctx, auth.BudgetID(), transactionIDs)
}

func (c *transactionContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.TransactionWithRelations, error) {
	return c.ds().TransactionRepository().GetForBudget(ctx, auth.BudgetID())
}

func (c *transactionContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.TransactionWithRelations, error) {
	transaction, err := c.ds().TransactionRepository().Get(ctx, id)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	currentAccount, err := c.ds().AccountRepository().Get(ctx, transaction.AccountID)
	if err != nil {
		return beans.TransactionWithRelations{}, fmt.Errorf("could not find related account: %w", err)
	}
	if currentAccount.BudgetID != auth.BudgetID() {
		return beans.TransactionWithRelations{}, beans.NewError(beans.ENOTFOUND, "transaction not found")
	}

	fullTransaction := beans.TransactionWithRelations{
		Transaction: transaction,
		Account:     beans.RelatedAccount{ID: currentAccount.ID, Name: currentAccount.Name},
	}

	if !transaction.CategoryID.Empty() {
		category, err := c.ds().CategoryRepository().GetSingleForBudget(ctx, transaction.CategoryID, auth.BudgetID())
		if err != nil {
			return beans.TransactionWithRelations{}, fmt.Errorf("could not find related category %s for transaction %s", transaction.CategoryID, id)
		}

		fullTransaction.Category = beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name})
	}

	if !transaction.PayeeID.Empty() {
		payee, err := c.ds().PayeeRepository().Get(ctx, transaction.PayeeID)
		if err != nil {
			return beans.TransactionWithRelations{}, fmt.Errorf("could not find related payee %s for transaction %s", transaction.PayeeID, id)
		}

		fullTransaction.Payee = beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name})
	}

	return fullTransaction, nil
}

func (c *transactionContract) validatePayee(ctx context.Context, auth *beans.BudgetAuthContext, payeeID beans.ID) error {
	if !payeeID.Empty() {
		payee, err := c.ds().PayeeRepository().Get(ctx, payeeID)

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
