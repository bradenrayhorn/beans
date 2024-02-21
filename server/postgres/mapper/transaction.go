package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
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

	return beans.Transaction{
		ID:         id,
		AccountID:  accountID,
		CategoryID: categoryID,
		PayeeID:    payeeID,
		Amount:     amount,
		Date:       PgToDate(d.Date),
		Notes:      beans.TransactionNotes{NullString: PgToNullString(d.Notes)},
	}, nil
}

func GetTransactionsForBudgetRow(d db.GetTransactionsForBudgetRow) (beans.TransactionWithRelations, error) {
	transaction, err := Transaction(d.Transaction)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	categoryName := PgToNullString(d.CategoryName)
	payeeName := PgToNullString(d.PayeeName)

	transactionWithRelations := beans.TransactionWithRelations{
		Transaction: transaction,
		Account: beans.RelatedAccount{
			ID:   transaction.AccountID,
			Name: beans.Name(d.AccountName),
		},
	}

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
