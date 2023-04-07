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
		PayeeID:    idToNullString(transaction.PayeeID),
		Date:       transaction.Date.Time,
		Amount:     amountToNumeric(transaction.Amount),
		Notes:      transaction.Notes.SQLNullString(),
	})
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *beans.Transaction) error {
	return r.db.UpdateTransaction(ctx, db.UpdateTransactionParams{
		ID:         transaction.ID.String(),
		AccountID:  transaction.AccountID.String(),
		CategoryID: idToNullString(transaction.CategoryID),
		PayeeID:    idToNullString(transaction.PayeeID),
		Date:       transaction.Date.Time,
		Amount:     amountToNumeric(transaction.Amount),
		Notes:      transaction.Notes.SQLNullString(),
	})
}

func (r *TransactionRepository) Get(ctx context.Context, id beans.ID) (*beans.Transaction, error) {
	t, err := r.db.GetTransaction(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	budgetID, err := beans.BeansIDFromString(t.BudgetID)
	if err != nil {
		return nil, err
	}
	accountID, err := beans.BeansIDFromString(t.AccountID)
	if err != nil {
		return nil, err
	}
	amount, err := numericToAmount(t.Amount)
	if err != nil {
		return nil, err
	}
	categoryID, err := nullStringToID(t.CategoryID)
	if err != nil {
		return nil, err
	}
	payeeID, err := nullStringToID(t.PayeeID)
	if err != nil {
		return nil, err
	}

	return &beans.Transaction{
		ID:         id,
		AccountID:  accountID,
		CategoryID: categoryID,
		PayeeID:    payeeID,
		Amount:     amount,
		Date:       beans.NewDate(t.Date),
		Notes:      beans.TransactionNotes{NullString: beans.NullStringFromSQL(t.Notes)},
		Account: &beans.Account{
			ID:       accountID,
			Name:     beans.Name(t.AccountName),
			BudgetID: budgetID,
		},
	}, nil
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
		payeeID, err := nullStringToID(t.PayeeID)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, &beans.Transaction{
			ID:         id,
			AccountID:  accountID,
			CategoryID: categoryID,
			PayeeID:    payeeID,
			Amount:     amount,
			Date:       beans.NewDate(t.Date),
			Notes:      beans.TransactionNotes{NullString: beans.NullStringFromSQL(t.Notes)},
			Account: &beans.Account{
				ID:       accountID,
				Name:     beans.Name(t.AccountName),
				BudgetID: budgetID,
			},
			CategoryName: beans.NullStringFromSQL(t.CategoryName),
			PayeeName:    beans.NullStringFromSQL(t.PayeeName),
		})
	}

	return transactions, nil
}

func (r *TransactionRepository) GetIncomeBetween(ctx context.Context, budgetID beans.ID, begin beans.Date, end beans.Date) (beans.Amount, error) {
	res, err := r.db.GetIncomeBetween(ctx, db.GetIncomeBetweenParams{
		BudgetID:  budgetID.String(),
		BeginDate: begin.Time,
		EndDate:   end.Time,
	})
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	amount, err := numericToAmount(res)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	return amount.OrZero(), nil
}
