package mapper

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/jackc/pgx/v5/pgtype"
)

func IDToPg(id beans.ID) pgtype.Text {
	return pgtype.Text{
		String: id.String(),
		Valid:  !id.Empty(),
	}
}

func PgToID(pg pgtype.Text) (beans.ID, error) {
	if !pg.Valid {
		return beans.IDFromString("")
	} else {
		return beans.IDFromString(pg.String)
	}
}

func NullStringToPg(s beans.NullString) pgtype.Text {
	return pgtype.Text{
		String: s.String(),
		Valid:  !s.Empty(),
	}
}

func PgToNullString(pg pgtype.Text) beans.NullString {
	if pg.Valid {
		return beans.NewNullString(pg.String)
	} else {
		return beans.NewNullString("")
	}
}
