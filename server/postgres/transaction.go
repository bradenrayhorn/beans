package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type TransactionRepository struct {
	repository
}

func NewTransactionRepository(pool *DbPool) *TransactionRepository {
	return &TransactionRepository{repository{pool}}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *beans.Transaction) error {
	return r.DB(nil).CreateTransaction(ctx, db.CreateTransactionParams{
		ID:         transaction.ID.String(),
		AccountID:  transaction.AccountID.String(),
		CategoryID: mapper.IDToPg(transaction.CategoryID),
		PayeeID:    mapper.IDToPg(transaction.PayeeID),
		Date:       mapper.DateToPg(transaction.Date),
		Amount:     mapper.AmountToNumeric(transaction.Amount),
		Notes:      mapper.NullStringToPg(transaction.Notes.NullString),
	})
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *beans.Transaction) error {
	return r.DB(nil).UpdateTransaction(ctx, db.UpdateTransactionParams{
		ID:         transaction.ID.String(),
		AccountID:  transaction.AccountID.String(),
		CategoryID: mapper.IDToPg(transaction.CategoryID),
		PayeeID:    mapper.IDToPg(transaction.PayeeID),
		Date:       mapper.DateToPg(transaction.Date),
		Amount:     mapper.AmountToNumeric(transaction.Amount),
		Notes:      mapper.NullStringToPg(transaction.Notes.NullString),
	})
}

func (r *TransactionRepository) Delete(ctx context.Context, budgetID beans.ID, transactionIDs []beans.ID) error {
	return r.DB(nil).DeleteTransactions(ctx, db.DeleteTransactionsParams{
		BudgetID: budgetID.String(),
		Ids: mapper.MapSliceNoErr(transactionIDs, func(id beans.ID) string {
			return id.String()
		}),
	})
}

func (r *TransactionRepository) Get(ctx context.Context, id beans.ID) (*beans.Transaction, error) {
	t, err := r.DB(nil).GetTransaction(ctx, id.String())
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
	amount, err := mapper.NumericToAmount(t.Amount)
	if err != nil {
		return nil, err
	}
	categoryID, err := mapper.PgToID(t.CategoryID)
	if err != nil {
		return nil, err
	}
	payeeID, err := mapper.PgToID(t.PayeeID)
	if err != nil {
		return nil, err
	}

	return &beans.Transaction{
		ID:         id,
		AccountID:  accountID,
		CategoryID: categoryID,
		PayeeID:    payeeID,
		Amount:     amount,
		Date:       mapper.PgToDate(t.Date),
		Notes:      beans.TransactionNotes{NullString: mapper.PgToNullString(t.Notes)},
		Account: &beans.Account{
			ID:       accountID,
			Name:     beans.Name(t.AccountName),
			BudgetID: budgetID,
		},
	}, nil
}

func (r *TransactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Transaction, error) {
	transactions := []*beans.Transaction{}
	dbTransactions, err := r.DB(nil).GetTransactionsForBudget(ctx, budgetID.String())
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
		amount, err := mapper.NumericToAmount(t.Amount)
		if err != nil {
			return transactions, err
		}
		categoryID, err := mapper.PgToID(t.CategoryID)
		if err != nil {
			return transactions, err
		}
		payeeID, err := mapper.PgToID(t.PayeeID)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, &beans.Transaction{
			ID:         id,
			AccountID:  accountID,
			CategoryID: categoryID,
			PayeeID:    payeeID,
			Amount:     amount,
			Date:       mapper.PgToDate(t.Date),
			Notes:      beans.TransactionNotes{NullString: mapper.PgToNullString(t.Notes)},
			Account: &beans.Account{
				ID:       accountID,
				Name:     beans.Name(t.AccountName),
				BudgetID: budgetID,
			},
			CategoryName: mapper.PgToNullString(t.CategoryName),
			PayeeName:    mapper.PgToNullString(t.PayeeName),
		})
	}

	return transactions, nil
}

func (r *TransactionRepository) GetIncomeBetween(ctx context.Context, budgetID beans.ID, begin beans.Date, end beans.Date) (beans.Amount, error) {
	res, err := r.DB(nil).GetIncomeBetween(ctx, db.GetIncomeBetweenParams{
		BudgetID:  budgetID.String(),
		BeginDate: mapper.DateToPg(begin),
		EndDate:   mapper.DateToPg(end),
	})
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	amount, err := mapper.NumericToAmount(res)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	return amount.OrZero(), nil
}
