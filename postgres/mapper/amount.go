package mapper

import (
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/jackc/pgtype"
)

func NumericToAmount(n pgtype.Numeric) (beans.Amount, error) {
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
