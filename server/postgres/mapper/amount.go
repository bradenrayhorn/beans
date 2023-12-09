package mapper

import (
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToAmount(n pgtype.Numeric) (beans.Amount, error) {
	if !n.Valid {
		return beans.NewEmptyAmount(), nil
	}

	if n.NaN {
		return beans.Amount{}, errors.New("invalid amount")
	}

	return beans.NewAmountWithBigInt(n.Int, n.Exp), nil
}

func AmountToNumeric(a beans.Amount) pgtype.Numeric {
	valid := true
	if a.Empty() {
		valid = false
	}
	return pgtype.Numeric{
		Int:              a.Coefficient(),
		Exp:              a.Exponent(),
		Valid:            valid,
		InfinityModifier: pgtype.Finite,
	}
}
