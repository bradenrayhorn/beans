package postgres

import (
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/jackc/pgx/v5"
)

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return beans.WrapError(err, beans.ErrorNotFound)
	}
	return err
}
