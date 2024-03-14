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

	account, err := c.validateAccount(ctx, auth, data.AccountID, "Invalid Account ID")
	if err != nil {
		return beans.EmptyID(), err
	}

	if err = c.validateCategory(ctx, auth, account, data.CategoryID); err != nil {
		return beans.EmptyID(), err
	}

	if err = c.validatePayee(ctx, auth, data.PayeeID); err != nil {
		return beans.EmptyID(), err
	}

	transaction := beans.Transaction{
		ID:         beans.NewID(),
		AccountID:  data.AccountID,
		CategoryID: data.CategoryID,
		PayeeID:    data.PayeeID,
		Amount:     data.Amount,
		Date:       data.Date,
		Notes:      data.Notes,
	}

	transactions := []beans.Transaction{transaction}

	// check if initiating transfer
	if !data.TransferAccountID.Empty() {

		if _, err := c.validateAccount(ctx, auth, data.TransferAccountID, "Invalid Transfer Account"); err != nil {
			return beans.EmptyID(), err
		}

		// copy transaction
		transactionB := beans.Transaction{
			ID:         beans.NewID(),
			AccountID:  data.TransferAccountID,
			Amount:     beans.Arithmetic.Negate(data.Amount),
			Date:       data.Date,
			Notes:      data.Notes,
			TransferID: transaction.ID,
		}

		// make transactionA a transfer
		transaction.TransferID = transactionB.ID

		// save
		transactions = []beans.Transaction{transaction, transactionB}
	}

	err = c.ds().TransactionRepository().Create(ctx, transactions)
	if err != nil {
		return beans.EmptyID(), err
	}

	return transaction.ID, nil
}

func (c *transactionContract) Update(ctx context.Context, auth *beans.BudgetAuthContext, data beans.TransactionUpdateParams) error {
	if err := data.ValidateAll(); err != nil {
		return err
	}

	transaction, err := c.ds().TransactionRepository().Get(ctx, auth.BudgetID(), data.ID)
	if err != nil {
		return err
	}

	account, err := c.validateAccount(ctx, auth, data.AccountID, "Invalid Account ID")
	if err != nil {
		return err
	}

	// cannot set payee or category on transfer
	if !transaction.TransferID.Empty() && (!data.CategoryID.Empty() || !data.PayeeID.Empty()) {
		return beans.NewError(beans.EINVALID, "cannot set a payee or category on transfer")
	}

	if err = c.validateCategory(ctx, auth, account, data.CategoryID); err != nil {
		return err
	}

	if err = c.validatePayee(ctx, auth, data.PayeeID); err != nil {
		return err
	}

	// update primary transaction
	transaction.AccountID = data.AccountID
	transaction.CategoryID = data.CategoryID
	transaction.PayeeID = data.PayeeID
	transaction.Amount = data.Amount
	transaction.Date = data.Date
	transaction.Notes = data.Notes

	updates := []beans.Transaction{transaction}

	// update transfer, if it exists
	if !transaction.TransferID.Empty() {
		transactionB, err := c.ds().TransactionRepository().Get(ctx, auth.BudgetID(), transaction.TransferID)
		if err != nil {
			return fmt.Errorf("could not get transfer: %w", err)
		}

		transactionB.Amount = beans.Arithmetic.Negate(data.Amount)
		transactionB.Date = data.Date
		transactionB.Notes = data.Notes

		updates = append(updates, transactionB)
	}

	if err := c.ds().TransactionRepository().Update(ctx, updates); err != nil {
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
	return c.ds().TransactionRepository().GetWithRelations(ctx, auth.BudgetID(), id)
}

func (c *transactionContract) validatePayee(ctx context.Context, auth *beans.BudgetAuthContext, payeeID beans.ID) error {
	if !payeeID.Empty() {
		if _, err := c.ds().PayeeRepository().Get(ctx, auth.BudgetID(), payeeID); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.NewError(beans.EINVALID, "Invalid Payee ID")
			}

			return err
		}

	}

	return nil
}

func (c *transactionContract) validateCategory(ctx context.Context, auth *beans.BudgetAuthContext, account beans.Account, categoryID beans.ID) error {
	if !categoryID.Empty() {
		if account.OffBudget {
			return beans.NewError(beans.EINVALID, "Cannot assign category with off-budget account")
		}

		if _, err := c.ds().CategoryRepository().GetSingleForBudget(ctx, categoryID, auth.BudgetID()); err != nil {
			if errors.Is(err, beans.ErrorNotFound) {
				return beans.NewError(beans.EINVALID, "Invalid Category ID")
			}

			return err
		}
	}

	return nil
}

func (c *transactionContract) validateAccount(ctx context.Context, auth *beans.BudgetAuthContext, accountID beans.ID, msg string) (beans.Account, error) {
	account, err := c.ds().AccountRepository().Get(ctx, auth.BudgetID(), accountID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.Account{}, beans.NewError(beans.EINVALID, msg)
		}

		return beans.Account{}, err
	}

	return account, nil
}
