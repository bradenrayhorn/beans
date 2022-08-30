package beans_test

import (
	"encoding/json"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanCreateAmount(t *testing.T) {
	amount := beans.NewAmount(55, -1)

	assert.Equal(t, int64(55), amount.Coefficient().Int64())
	assert.Equal(t, int32(-1), amount.Exponent())
}

func TestCanCreateNegativeAmount(t *testing.T) {
	amount := beans.NewAmount(-55, -1)

	assert.Equal(t, int64(-55), amount.Coefficient().Int64())
	assert.Equal(t, int32(-1), amount.Exponent())
}

// empty tests

func TestNewAmountIsNotEmpty(t *testing.T) {
	amount := beans.NewAmount(55, -1)

	assert.False(t, amount.Empty())
}

func TestUnmarshaledAmountIsNotEmpty(t *testing.T) {
	var amount beans.Amount
	err := json.Unmarshal([]byte("55"), &amount)
	require.Nil(t, err)

	assert.False(t, amount.Empty())
}

func TestBlankAmountIsEmpty(t *testing.T) {
	var amount beans.Amount
	assert.True(t, amount.Empty())
}
