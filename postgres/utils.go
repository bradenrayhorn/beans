package postgres

import (
	"database/sql"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func amountToNumeric(a beans.Amount) pgtype.Numeric {
	status := pgtype.Present
	if a.Empty() {
		status = pgtype.Null
	}
	return pgtype.Numeric{
		Int:              a.Coefficient(),
		Exp:              a.Exponent(),
		Status:           status,
		InfinityModifier: pgtype.None,
	}
}

func numericToAmount(n pgtype.Numeric) (beans.Amount, error) {
	if n.Status == pgtype.Null {
		return beans.NewEmptyAmount(), nil
	}

	if n.Status != pgtype.Present {
		return beans.Amount{}, errors.New("invalid amount")
	}
	if n.NaN {
		return beans.Amount{}, errors.New("invalid amount")
	}

	return beans.NewAmountWithBigInt(n.Int, n.Exp), nil
}

func idToNullString(id beans.ID) sql.NullString {
	if id.Empty() {
		return sql.NullString{String: "", Valid: false}
	} else {
		return sql.NullString{String: id.String(), Valid: true}
	}
}

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return beans.WrapError(err, beans.ErrorNotFound)
	}
	return err
}
