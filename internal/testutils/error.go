package testutils

import (
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertError(t testing.TB, err error, expected string) {
	require.NotNil(t, err)
	_, msg := err.(beans.Error).BeansError()
	assert.Equal(t, expected, msg)
}

func AssertErrorCode(t testing.TB, err error, expected string) {
	require.NotNil(t, err)
	code, _ := err.(beans.Error).BeansError()
	assert.Equal(t, expected, code)
}

func NewDate(t testing.TB, date string) beans.Date {
	time, err := time.Parse("2006-01-02", date)
	require.Nil(t, err)
	return beans.NewDate(time)
}
