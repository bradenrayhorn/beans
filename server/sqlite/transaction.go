package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type TransactionRepository struct{ repository }

var _ beans.TransactionRepository = (*TransactionRepository)(nil)

func (r *TransactionRepository) Create(ctx context.Context, transactions []beans.Transaction) error {
	q := squirrel.
		Insert("transactions").
		Columns("id", "account_id", "category_id", "payee_id", "amount", "date", "notes", "transfer_id", "split_id", "is_split")

	for _, t := range transactions {
		amount, err := serializeAmount(t.Amount)
		if err != nil {
			return err
		}

		q = q.Values(
			t.ID.String(),
			serializeID(t.AccountID),
			serializeID(t.CategoryID),
			serializeID(t.PayeeID),
			amount,
			serializeDate(t.Date),
			serializeNullString(t.Notes.NullString),
			serializeID(t.TransferID),
			serializeID(t.SplitID),
			t.IsSplit,
		)
	}

	sql, params, err := q.ToSql()
	if err != nil {
		return err
	}

	return db[any](r.pool).executeWithArgs(ctx, sql, params)
}

const updateTransactionSQL = `
UPDATE transactions
	SET account_id=:accountID, category_id=:categoryID, payee_id=:payeeID, date=:date, amount=:amount, notes=:notes, is_split=:isSplit
	WHERE id=:id
`

func (r *TransactionRepository) Update(ctx context.Context, transactions []beans.Transaction) error {
	txm := &txManager{r.pool}
	return beans.ExecTxNil(ctx, txm, func(tx beans.Tx) error {
		for _, t := range transactions {
			amount, err := serializeAmount(t.Amount)
			if err != nil {
				return err
			}

			err = db[any](r.pool).
				inTx(tx).
				execute(ctx, updateTransactionSQL, map[string]any{
					":id":         t.ID.String(),
					":accountID":  t.AccountID.String(),
					":categoryID": serializeID(t.CategoryID),
					":payeeID":    serializeID(t.PayeeID),
					":amount":     amount,
					":date":       serializeDate(t.Date),
					":notes":      serializeNullString(t.Notes.NullString),
					":isSplit":    t.IsSplit,
				})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TransactionRepository) Delete(ctx context.Context, budgetID beans.ID, transactionIDs []beans.ID) error {
	// get transaction ids to delete
	sql, params, err := squirrel.
		Select("transactions.id").
		From("transactions").
		Join("accounts ON accounts.id = transactions.account_id AND accounts.budget_id = ?", budgetID.String()).
		Where(squirrel.And{
			squirrel.Expr("transactions.split_id IS NULL"),
			squirrel.Eq{"transactions.id": transactionIDs},
		}).
		ToSql()
	if err != nil {
		return err
	}

	ids, err := db[string](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (string, error) { return stmt.GetText("id"), nil }).
		manyWithArgs(ctx, sql, params)
	if err != nil {
		return err
	}

	// delete transactions
	sql, params, err = squirrel.
		Delete("transactions").
		From("transactions").
		Where(squirrel.Eq{"id": ids}).
		ToSql()

	return db[any](r.pool).
		executeWithArgs(ctx, sql, params)
}

const transactionGetSQL = `
SELECT transactions.* FROM transactions
JOIN accounts ON accounts.id = transactions.account_id
	AND accounts.budget_id = :budgetID
WHERE transactions.id = :id
`

func (r *TransactionRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Transaction, error) {
	return db[beans.Transaction](r.pool).
		mapWith(mapTransaction).
		one(ctx, transactionGetSQL, map[string]any{
			":budgetID": budgetID.String(),
			":id":       id.String(),
		})
}

func (r *TransactionRepository) GetWithRelations(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.TransactionWithRelations, error) {
	q := getTransactionWithRelationshipsQuery(budgetID.String()).
		Where("transactions.id = ?", id)
	sql, args, err := q.ToSql()
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	return db[beans.TransactionWithRelations](r.pool).
		mapWith(mapTransactionWithRelations).
		oneWithArgs(ctx, sql, args)
}

func (r *TransactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.TransactionWithRelations, error) {
	q := getTransactionWithRelationshipsQuery(budgetID.String()).
		Where("transactions.split_id IS NULL")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	return db[beans.TransactionWithRelations](r.pool).
		mapWith(mapTransactionWithRelations).
		manyWithArgs(ctx, sql, args)
}

func (r *TransactionRepository) GetSplits(ctx context.Context, budgetID beans.ID, transactionID beans.ID) ([]beans.TransactionAsSplit, error) {
	q := getTransactionWithRelationshipsQuery(budgetID.String()).
		Where("transactions.split_id = ?", transactionID.String())
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	return db[beans.TransactionAsSplit](r.pool).
		mapWith(mapTransactionAsSplit).
		manyWithArgs(ctx, sql, args)
}

type getActivityByCategoryRow struct {
	ID       beans.ID
	Activity beans.Amount
}

func (r *TransactionRepository) GetActivityByCategory(ctx context.Context, budgetID beans.ID, from beans.Date, to beans.Date) (map[beans.ID]beans.Amount, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := psql.
		Select("categories.id", "sum(transactions.amount) as activity").
		From("transactions").
		Join("categories ON transactions.category_id = categories.id").
		Join("accounts ON transactions.account_id = accounts.id AND accounts.budget_id = ?", budgetID.String()).
		GroupBy("categories.id")

	if !from.Empty() {
		q = q.Where("transactions.date >= ?", serializeDate(from))
	}
	if !to.Empty() {
		q = q.Where("transactions.date <= ?", serializeDate(to))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db[getActivityByCategoryRow](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (getActivityByCategoryRow, error) {
			id, err := mapID(stmt, "id")
			if err != nil {
				return getActivityByCategoryRow{}, err
			}
			return getActivityByCategoryRow{ID: id, Activity: mapAmount(stmt, "activity")}, nil
		}).
		manyWithArgs(ctx, sql, args)
	if err != nil {
		return nil, err
	}

	activityByCategory := make(map[beans.ID]beans.Amount)
	for _, v := range rows {
		activityByCategory[v.ID] = v.Activity
	}

	return activityByCategory, nil
}

const transactionGetIncomeBetweenSQL = `
SELECT sum(transactions.amount) as income
FROM transactions
JOIN categories
  ON categories.id = transactions.category_id
JOIN category_groups
  ON category_groups.id = categories.group_id
  AND category_groups.is_income = true
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = :budgetID
WHERE
	transactions.date <= :end
	AND transactions.date >= :begin
`

func (r *TransactionRepository) GetIncomeBetween(ctx context.Context, budgetID beans.ID, begin beans.Date, end beans.Date) (beans.Amount, error) {
	return db[beans.Amount](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (beans.Amount, error) { return mapAmount(stmt, "income"), nil }).
		one(ctx, transactionGetIncomeBetweenSQL, map[string]any{
			":budgetID": budgetID.String(),
			":begin":    serializeDate(begin),
			":end":      serializeDate(end),
		})
}

// big queries

func getTransactionWithRelationshipsQuery(budgetID string) squirrel.SelectBuilder {
	return squirrel.
		Select(
			"transactions.*",
			"accounts.name as account_name",
			"categories.name as category_name",
			"payees.name as payee_name",
			"accounts.off_budget as account_off_budget",
			"transfer_account.id as transfer_account_id",
			"transfer_account.name as transfer_account_name",
			"transfer_account.off_budget as transfer_account_off_budget",
		).
		From("transactions").
		Join("accounts ON transactions.account_id = accounts.id AND accounts.budget_id = ?", budgetID).
		LeftJoin("categories ON categories.id = transactions.category_id").
		LeftJoin("payees ON payees.id = transactions.payee_id").
		LeftJoin("transactions transfer ON transfer.id = transactions.transfer_id").
		LeftJoin("accounts transfer_account ON transfer.account_id = transfer_account.id").
		OrderBy("transactions.date DESC")
}

// mappers

func mapTransaction(stmt *sqlite.Stmt) (beans.Transaction, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Transaction{}, err
	}
	accountID, err := mapID(stmt, "account_id")
	if err != nil {
		return beans.Transaction{}, err
	}
	categoryID, err := mapID(stmt, "category_id")
	if err != nil {
		return beans.Transaction{}, err
	}
	payeeID, err := mapID(stmt, "payee_id")
	if err != nil {
		return beans.Transaction{}, err
	}
	transferID, err := mapID(stmt, "transfer_id")
	if err != nil {
		return beans.Transaction{}, err
	}
	splitID, err := mapID(stmt, "split_id")
	if err != nil {
		return beans.Transaction{}, err
	}
	date, err := mapDate(stmt, "date")
	if err != nil {
		return beans.Transaction{}, err
	}

	return beans.Transaction{
		ID:         id,
		AccountID:  accountID,
		CategoryID: categoryID,
		PayeeID:    payeeID,

		Amount: mapAmount(stmt, "amount"),
		Date:   date,
		Notes:  beans.TransactionNotes{NullString: mapNullString(stmt, "notes")},

		TransferID: transferID,
		SplitID:    splitID,
		IsSplit:    stmt.GetBool("is_split"),
	}, nil
}

func mapTransactionWithRelations(stmt *sqlite.Stmt) (beans.TransactionWithRelations, error) {
	transaction, err := mapTransaction(stmt)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	categoryName := mapNullString(stmt, "category_name")
	payeeName := mapNullString(stmt, "payee_name")

	transactionWithRelations := beans.TransactionWithRelations{
		ID:     transaction.ID,
		Amount: transaction.Amount,
		Date:   transaction.Date,
		Notes:  transaction.Notes,
		Account: beans.RelatedAccount{
			ID:        transaction.AccountID,
			Name:      beans.Name(stmt.GetText("account_name")),
			OffBudget: stmt.GetBool("account_off_budget"),
		},
	}

	if !transaction.TransferID.Empty() {
		transferAccountID, err := mapID(stmt, "transfer_account_id")
		if err != nil {
			return beans.TransactionWithRelations{}, err
		}

		transactionWithRelations.TransferAccount = beans.OptionalWrap(beans.RelatedAccount{
			ID:        transferAccountID,
			Name:      beans.Name(stmt.GetText("transfer_account_name")),
			OffBudget: stmt.GetBool("transfer_account_off_budget"),
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

func mapTransactionAsSplit(stmt *sqlite.Stmt) (beans.TransactionAsSplit, error) {
	transaction, err := mapTransaction(stmt)
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
				Name: beans.Name(stmt.GetText("category_name")),
			},
		},
	}, nil
}
