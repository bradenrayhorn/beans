package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/jackc/pgx/v5/pgtype"
)

func DateToPg(date beans.Date) pgtype.Date {
	return pgtype.Date{
		Time:             date.Time,
		Valid:            !date.Empty(),
		InfinityModifier: pgtype.Finite,
	}
}

func PgToDate(pg pgtype.Date) beans.Date {
	if pg.Valid {
		return beans.NewDate(pg.Time)
	} else {
		return beans.Date{}
	}
}

func MonthDateToPg(date beans.MonthDate) pgtype.Date {
	return pgtype.Date{
		Time:             date.FirstDay().Time,
		Valid:            !date.FirstDay().Empty(),
		InfinityModifier: pgtype.Finite,
	}
}

func PgToMonthDate(pg pgtype.Date) beans.MonthDate {
	if pg.Valid {
		return beans.NewMonthDate(beans.NewDate(pg.Time))
	} else {
		return beans.NewMonthDate(beans.Date{})
	}
}
