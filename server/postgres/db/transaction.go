package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// models

type Transaction struct {
	ID         string           `db:"id"`
	AccountID  string           `db:"account_id"`
	PayeeID    pgtype.Text      `db:"payee_id"`
	CategoryID pgtype.Text      `db:"category_id"`
	TransferID pgtype.Text      `db:"transfer_id"`
	SplitID    pgtype.Text      `db:"split_id"`
	IsSplit    bool             `db:"is_split"`
	Date       pgtype.Date      `db:"date"`
	Amount     pgtype.Numeric   `db:"amount"`
	Notes      pgtype.Text      `db:"notes"`
	CreatedAt  pgtype.Timestamp `db:"created_at"`
}

// CreateTransaction

type CreateTransactionParams struct {
	ID         string
	AccountID  string
	PayeeID    pgtype.Text
	CategoryID pgtype.Text
	Date       pgtype.Date
	Amount     pgtype.Numeric
	Notes      pgtype.Text
	TransferID pgtype.Text
	SplitID    pgtype.Text
	IsSplit    bool
}

func (e *Executor) CreateTransactions(ctx context.Context, params []CreateTransactionParams) error {
	_, err := e.db.CopyFrom(
		ctx,
		pgx.Identifier{"transactions"},
		[]string{"id", "account_id", "payee_id", "category_id", "date", "amount", "notes", "transfer_id", "split_id", "is_split"},
		pgx.CopyFromSlice(len(params), func(i int) ([]any, error) {
			return []any{
				params[i].ID,
				params[i].AccountID,
				params[i].PayeeID,
				params[i].CategoryID,
				params[i].Date,
				params[i].Amount,
				params[i].Notes,
				params[i].TransferID,
				params[i].SplitID,
				params[i].IsSplit,
			}, nil
		}),
	)
	return err
}

// UpdateTransaction

const updateTransactionSQL = `
UPDATE transactions
  SET account_id=$1, category_id=$2, payee_id=$3, date=$4, amount=$5, notes=$6, is_split=$7
  WHERE id=$8;
`

type UpdateTransactionParams struct {
	AccountID  string
	CategoryID pgtype.Text
	PayeeID    pgtype.Text
	Date       pgtype.Date
	Amount     pgtype.Numeric
	Notes      pgtype.Text
	IsSplit    bool
	ID         string
}

func (e *Executor) UpdateTransaction(ctx context.Context, params UpdateTransactionParams) error {
	_, err := e.db.Exec(ctx, updateTransactionSQL,
		params.AccountID,
		params.CategoryID,
		params.PayeeID,
		params.Date,
		params.Amount,
		params.Notes,
		params.IsSplit,
		params.ID,
	)
	return err
}

// DeleteTransactions

const deleteTransactionsSQL = `
DELETE FROM transactions
  USING accounts
  WHERE
    accounts.id = transactions.account_id
    AND accounts.budget_id=$1
    AND transactions.id = ANY($2)
	AND transactions.split_id IS NULL
`

type DeleteTransactionsParams struct {
	BudgetID string
	IDs      []string
}

func (e *Executor) DeleteTransactions(ctx context.Context, params DeleteTransactionsParams) error {
	_, err := e.db.Exec(ctx, deleteTransactionsSQL,
		params.BudgetID,
		params.IDs,
	)
	return err
}

// GetTransaction

const getTransactionSQL = `
SELECT transactions.*
  FROM transactions
  JOIN accounts
    ON accounts.id = transactions.account_id
    AND accounts.budget_id = $1
  WHERE transactions.id = $2;
`

type GetTransactionParams struct {
	ID       string
	BudgetID string
}

func (e *Executor) GetTransaction(ctx context.Context, params GetTransactionParams) (Transaction, error) {
	row, _ := e.db.Query(ctx, getTransactionSQL, params.BudgetID, params.ID)
	return pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[Transaction])
}

// GetIncomeBetween

const getIncomeBetweenSQL = `
SELECT sum(transactions.amount)::numeric
FROM transactions
JOIN categories
  ON categories.id = transactions.category_id
JOIN category_groups
  ON category_groups.id = categories.group_id
  AND category_groups.is_income = true
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1
WHERE
  transactions.date <= $2
  AND transactions.date >= $3
`

type GetIncomeBetweenParams struct {
	BudgetID  string
	BeginDate pgtype.Date
	EndDate   pgtype.Date
}

func (e *Executor) GetIncomeBetween(ctx context.Context, params GetIncomeBetweenParams) (pgtype.Numeric, error) {
	row, _ := e.db.Query(ctx, getIncomeBetweenSQL,
		params.BudgetID,
		params.EndDate,
		params.BeginDate,
	)
	return pgx.CollectExactlyOneRow(row, pgx.RowTo[pgtype.Numeric])
}

// GetActivityByCategory

type GetActivityByCategoryParams struct {
	BudgetID string
	FromDate pgtype.Date
	ToDate   pgtype.Date
}

type GetActivityByCategoryRow struct {
	ID       string
	Activity pgtype.Numeric
}

func (e *Executor) GetActivityByCategory(ctx context.Context, params GetActivityByCategoryParams) ([]GetActivityByCategoryRow, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := psql.
		Select("categories.id", "sum(transactions.amount)::numeric as activity").
		From("transactions").
		Join("categories ON transactions.category_id = categories.id").
		Join("accounts ON transactions.account_id = accounts.id AND accounts.budget_id = ?", params.BudgetID).
		GroupBy("categories.id")

	if params.FromDate.Valid {
		q = q.Where("transactions.date >= ?", params.FromDate)
	}
	if params.ToDate.Valid {
		q = q.Where("transactions.date <= ?", params.ToDate)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, _ := e.db.Query(ctx, sql, args...)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[GetActivityByCategoryRow])
}

// GetTransaction(s)ForBudget

type TransactionWithRelationships struct {
	Transaction
	AccountName              string
	CategoryName             pgtype.Text
	PayeeName                pgtype.Text
	AccountOffBudget         bool
	TransferAccountID        pgtype.Text
	TransferAccountName      pgtype.Text
	TransferAccountOffBudget pgtype.Bool
}

type GetTransactionWithRelationshipsParams struct {
	ID       string
	BudgetID string
}

type GetTransactionsWithRelationshipsParams struct {
	BudgetID string
}

type GetSplitsParams struct {
	TransactionID string
	BudgetID      string
}

func getTransactionWithRelationshipsQuery(budgetID string) squirrel.SelectBuilder {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return psql.
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

func (e *Executor) GetTransactionWithRelationships(ctx context.Context, params GetTransactionWithRelationshipsParams) (TransactionWithRelationships, error) {
	q := getTransactionWithRelationshipsQuery(params.BudgetID).
		Where("transactions.id = ?", params.ID)
	sql, args, err := q.ToSql()
	if err != nil {
		return TransactionWithRelationships{}, err
	}
	rows, _ := e.db.Query(ctx, sql, args...)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[TransactionWithRelationships])
}

func (e *Executor) GetTransactionsWithRelationships(ctx context.Context, params GetTransactionsWithRelationshipsParams) ([]TransactionWithRelationships, error) {
	q := getTransactionWithRelationshipsQuery(params.BudgetID).
		Where("transactions.split_id IS NULL")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, _ := e.db.Query(ctx, sql, args...)
	return pgx.CollectRows(rows, pgx.RowToStructByName[TransactionWithRelationships])
}

func (e *Executor) GetSplits(ctx context.Context, params GetSplitsParams) ([]TransactionWithRelationships, error) {
	q := getTransactionWithRelationshipsQuery(params.BudgetID).
		Where("transactions.split_id = ?", params.TransactionID)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, _ := e.db.Query(ctx, sql, args...)
	return pgx.CollectRows(rows, pgx.RowToStructByName[TransactionWithRelationships])
}
