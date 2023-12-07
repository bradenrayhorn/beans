package testutils

import (
	"errors"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertError(t testing.TB, err error, expected string) {
	require.NotNil(t, err)
	var beansError beans.Error
	require.True(t, errors.As(err, &beansError))
	_, msg := beansError.BeansError()
	assert.Equal(t, expected, msg)
}

func AssertErrorCode(t testing.TB, err error, expected string) {
	require.NotNil(t, err)
	var beansError beans.Error
	require.True(t, errors.As(err, &beansError))
	code, _ := beansError.BeansError()
	assert.Equal(t, expected, code)
}

func AssertErrorAndCode(t testing.TB, err error, code string, msg string) {
	AssertError(t, err, msg)
	AssertErrorCode(t, err, code)
}
