package mapper

import (
	"fmt"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/postgres/db"
)

func Transaction(d db.Transaction) (beans.Transaction, error) {
	id, err := beans.IDFromString(d.ID)
	if err != nil {
		return beans.Transaction{}, err
	}

	accountID, err := beans.IDFromString(d.AccountID)
	if err != nil {
		return beans.Transaction{}, err
	}
	amount, err := NumericToAmount(d.Amount)
	if err != nil {
		return beans.Transaction{}, err
	}
	categoryID, err := PgToID(d.CategoryID)
	if err != nil {
		return beans.Transaction{}, err
	}
	payeeID, err := PgToID(d.PayeeID)
	if err != nil {
		return beans.Transaction{}, err
	}
	transferID, err := PgToID(d.TransferID)
	if err != nil {
		return beans.Transaction{}, err
	}
	splitID, err := PgToID(d.SplitID)
	if err != nil {
		return beans.Transaction{}, err
	}

	return beans.Transaction{
		ID:         id,
		AccountID:  accountID,
		CategoryID: categoryID,
		PayeeID:    payeeID,
		Amount:     amount,
		Date:       PgToDate(d.Date),
		Notes:      beans.TransactionNotes{NullString: PgToNullString(d.Notes)},
		TransferID: transferID,
		SplitID:    splitID,
		IsSplit:    d.IsSplit,
	}, nil
}

func GetTransactionsForBudgetRow(d db.TransactionWithRelationships) (beans.TransactionWithRelations, error) {
	transaction, err := Transaction(d.Transaction)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	categoryName := PgToNullString(d.CategoryName)
	payeeName := PgToNullString(d.PayeeName)

	transactionWithRelations := beans.TransactionWithRelations{
		ID:     transaction.ID,
		Amount: transaction.Amount,
		Date:   transaction.Date,
		Notes:  transaction.Notes,
		Account: beans.RelatedAccount{
			ID:        transaction.AccountID,
			Name:      beans.Name(d.AccountName),
			OffBudget: d.AccountOffBudget,
		},
	}

	if !transaction.TransferID.Empty() {
		transferAccountID, err := PgToID(d.TransferAccountID)
		if err != nil {
			return beans.TransactionWithRelations{}, err
		}

		transactionWithRelations.TransferAccount = beans.OptionalWrap(beans.RelatedAccount{
			ID:        transferAccountID,
			Name:      beans.Name(PgToNullString(d.TransferAccountName).String()),
			OffBudget: d.TransferAccountOffBudget.Valid && d.TransferAccountOffBudget.Bool,
		})
	}

	transactionWithRelations.Variant = beans.GetTransactionVariant(
		transactionWithRelations.Account,
		transactionWithRelations.TransferAccount,
		transaction.IsSplit,
	)

	if !categoryName.Empty() {
		transactionWithRelations.Category = beans.OptionalWrap(beans.RelatedCategory{
			ID:   transaction.CategoryID,
			Name: beans.Name(categoryName.String()),
		})
	}

	if !payeeName.Empty() {
		transactionWithRelations.Payee = beans.OptionalWrap(beans.RelatedPayee{
			ID:   transaction.PayeeID,
			Name: beans.Name(payeeName.String()),
		})
	}

	return transactionWithRelations, nil
}

func Split(d db.TransactionWithRelationships) (beans.TransactionAsSplit, error) {
	transaction, err := Transaction(d.Transaction)
	if err != nil {
		return beans.TransactionAsSplit{}, err
	}

	if transaction.CategoryID.Empty() {
		return beans.TransactionAsSplit{}, fmt.Errorf("category null on split %s", transaction.ID)
	}

	return beans.TransactionAsSplit{
		Transaction: transaction,
		Split: beans.Split{
			ID:     transaction.ID,
			Amount: transaction.Amount,
			Notes:  transaction.Notes,
			Category: beans.RelatedCategory{
				ID:   transaction.CategoryID,
				Name: beans.Name(PgToNullString(d.CategoryName).String()),
			},
		},
	}, nil
}
