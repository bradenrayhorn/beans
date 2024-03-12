package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type TransactionRepository struct{ repository }

var _ beans.TransactionRepository = (*TransactionRepository)(nil)

func (r *TransactionRepository) Create(ctx context.Context, transaction beans.Transaction) error {
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

func (r *TransactionRepository) Update(ctx context.Context, transaction beans.Transaction) error {
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

func (r *TransactionRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Transaction, error) {
	t, err := r.DB(nil).GetTransaction(ctx, db.GetTransactionParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return beans.Transaction{}, mapPostgresError(err)
	}

	return mapper.Transaction(t)
}

func (r *TransactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.TransactionWithRelations, error) {
	res, err := r.DB(nil).GetTransactionsForBudget(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.GetTransactionsForBudgetRow)
}

func (r *TransactionRepository) GetActivityByCategory(ctx context.Context, budgetID beans.ID, from beans.Date, to beans.Date) (map[beans.ID]beans.Amount, error) {
	res, err := r.DB(nil).GetActivityByCategory(ctx, db.GetActivityByCategoryParams{
		BudgetID:       budgetID.String(),
		FromDate:       mapper.DateToPg(from),
		FilterFromDate: !from.Empty(),
		ToDate:         mapper.DateToPg(to),
		FilterToDate:   !to.Empty(),
	})
	if err != nil {
		return nil, err
	}

	activityByCategory := make(map[beans.ID]beans.Amount)
	for _, v := range res {
		id, err := beans.IDFromString(v.ID)
		if err != nil {
			return nil, err
		}
		amount, err := mapper.NumericToAmount(v.Activity)
		if err != nil {
			return nil, err
		}
		activityByCategory[id] = amount
	}

	return activityByCategory, nil
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
