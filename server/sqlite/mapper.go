package sqlite

import (
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

// id

func mapID(stmt *sqlite.Stmt, col string) (beans.ID, error) {
	id, err := beans.IDFromString(stmt.GetText(col))
	if err != nil {
		return beans.EmptyID(), err
	}
	return id, nil
}

func serializeID(id beans.ID) any {
	if id.Empty() {
		return nil
	} else {
		return id.String()
	}
}

// date

func mapDate(stmt *sqlite.Stmt, col string) (beans.Date, error) {
	if stmt.IsNull(col) {
		return beans.Date{}, nil
	} else {
		time, err := time.Parse("2006-01-02", stmt.GetText(col))
		if err != nil {
			return beans.Date{}, err
		}
		return beans.NewDate(time), nil
	}
}

func serializeDate(date beans.Date) any {
	if date.Empty() {
		return nil
	} else {
		return date.String()
	}
}

// null string

func mapNullString(stmt *sqlite.Stmt, col string) beans.NullString {
	if stmt.IsNull(col) {
		return beans.NullString{}
	} else {
		return beans.NewNullString(stmt.GetText(col))
	}
}

func serializeNullString(ns beans.NullString) any {
	if ns.Empty() {
		return nil
	} else {
		return ns.String()
	}
}

// amount

func serializeAmount(amount beans.Amount) (int64, error) {
	scaled := beans.NewAmountWithBigInt(amount.Coefficient(), amount.Exponent()+2)
	return scaled.AsInt64()
}

func mapAmount(stmt *sqlite.Stmt, col string) beans.Amount {
	return beans.NewAmount(stmt.GetInt64(col), -2).Normalize()
}
