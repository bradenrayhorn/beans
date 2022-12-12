package testutils

import (
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/require"
)

func NewDate(t testing.TB, date string) beans.Date {
	time, err := time.Parse("2006-01-02", date)
	require.Nil(t, err)
	return beans.NewDate(time)
}

func NewEmptyID() beans.ID {
	var id beans.ID
	return id
}
