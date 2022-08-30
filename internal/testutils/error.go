package testutils

import (
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertError(t testing.TB, err error, expected string) {
	require.NotNil(t, err)
	_, msg := err.(beans.Error).BeansError()
	assert.Equal(t, expected, msg)
}
