package postgres

import (
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func amountToNumeric(a beans.Amount) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              a.Coefficient(),
		Exp:              a.Exponent(),
		Status:           pgtype.Present,
		InfinityModifier: pgtype.None,
	}
}

func numericToAmount(n pgtype.Numeric) (beans.Amount, error) {
	if n.Status != pgtype.Present {
		return beans.Amount{}, errors.New("invalid amount")
	}
	if n.NaN {
		return beans.Amount{}, errors.New("invalid amount")
	}

	return beans.NewAmountWithBigInt(n.Int, n.Exp), nil
}

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return beans.WrapError(err, beans.ErrorNotFound)
	}
	return err
}
