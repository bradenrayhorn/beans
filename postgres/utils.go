package postgres

import (
	"github.com/bradenrayhorn/beans/beans"
	"github.com/jackc/pgtype"
)

func amountToNumeric(a beans.Amount) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              a.Coefficient(),
		Exp:              a.Exponent(),
		Status:           pgtype.Present,
		InfinityModifier: pgtype.None,
	}
}
