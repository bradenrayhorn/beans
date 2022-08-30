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

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return beans.WrapError(err, beans.ErrorNotFound)
	}
	return err
}
