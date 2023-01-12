package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository struct {
	db *db.Queries
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db.New(pool)}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *beans.Transaction) error {
	return r.db.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:         transaction.ID.String(),
		AccountID:  transaction.AccountID.String(),
		CategoryID: idToNullString(transaction.CategoryID),
		Date:       transaction.Date.Time,
		Amount:     amountToNumeric(transaction.Amount),
		Notes:      transaction.Notes.SQLNullString(),
	})
}

func (r *TransactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Transaction, error) {
	transactions := []*beans.Transaction{}
	dbTransactions, err := r.db.GetTransactionsForBudget(ctx, budgetID.String())
	if err != nil {
		return nil, nil
	}

	for _, t := range dbTransactions {
		id, err := beans.BeansIDFromString(t.ID)
		if err != nil {
			return transactions, err
		}
		accountID, err := beans.BeansIDFromString(t.AccountID)
		if err != nil {
			return transactions, err
		}
		amount, err := numericToAmount(t.Amount)
		if err != nil {
			return transactions, err
		}
		categoryID, err := nullStringToID(t.CategoryID)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, &beans.Transaction{
			ID:         id,
			AccountID:  accountID,
			CategoryID: categoryID,
			Amount:     amount,
			Date:       beans.NewDate(t.Date),
			Notes:      beans.TransactionNotes{NullString: beans.NullStringFromSQL(t.Notes)},
			Account: &beans.Account{
				ID:       accountID,
				Name:     beans.Name(t.AccountName),
				BudgetID: budgetID,
			},
			CategoryName: beans.NullStringFromSQL(t.CategoryName),
		})
	}

	return transactions, nil
}

func (r *TransactionRepository) GetIncomeBeforeOrOnDate(ctx context.Context, date beans.Date) (beans.Amount, error) {
	res, err := r.db.GetIncomeBeforeOrOnDate(ctx, date.Time)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	amount, err := numericToAmount(res)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	return amount, nil
}
