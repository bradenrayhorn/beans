package beans_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
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

func TestCanNormalizeAmount(t *testing.T) {
	// $1
	amount := beans.NewAmount(100, -2).Normalize()

	assert.Equal(t, int64(1), amount.Coefficient().Int64())
	assert.Equal(t, int32(0), amount.Exponent())

	// $500
	amount = beans.NewAmount(50000, -2).Normalize()

	assert.Equal(t, int64(5), amount.Coefficient().Int64())
	assert.Equal(t, int32(2), amount.Exponent())
}

func TestAmountJSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		var tests = []struct {
			name     string
			amount   beans.Amount
			expected string
		}{
			{"blank", beans.Amount{}, `null`},
			{"negative", beans.NewAmount(-5, 0), `"-5"`},
			{"positive", beans.NewAmount(5, 0), `"5"`},
			{"decimal", beans.NewAmount(55, -1), `"5.5"`},
			{"kind of big", beans.NewAmount(55, 3), `"55000"`},
			{"really big", beans.NewAmount(55, 10), `"550000000000"`},
			{"many decimals", beans.NewAmount(817401, -5), `"8.17401"`},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				r, err := json.Marshal(test.amount)
				require.Nil(t, err)
				assert.Equal(t, test.expected, string(r))
			})
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		var tests = []struct {
			name     string
			json     string
			expected string
		}{
			{"blank", `null`, ``},
			{"negative", `-50.12`, `-50.12`},
			{"positive", `60.76`, `60.76`},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var amount beans.Amount
				err := json.Unmarshal([]byte(test.json), &amount)
				require.Nil(t, err)
				assert.Equal(t, test.expected, amount.String())
			})
		}
	})
}

func TestNewAmountWithBigInt(t *testing.T) {
	t.Run("negative value", func(t *testing.T) {
		amount := beans.NewAmountWithBigInt(big.NewInt(-57), -1)
		assert.Equal(t, "-5.7", amount.String())
	})

	t.Run("positive value", func(t *testing.T) {
		amount := beans.NewAmountWithBigInt(big.NewInt(57), -1)
		assert.Equal(t, "5.7", amount.String())
	})
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

// non zero tests

func TestAmountNonZeroValidation(t *testing.T) {

	t.Run("filled in amount is not zero", func(t *testing.T) {
		amount := beans.NewAmount(55, -1)
		assert.Nil(t, beans.NonZero(amount).Validate())
	})

	t.Run("zero amount is zero", func(t *testing.T) {
		amount := beans.NewAmount(0, 1)
		assert.NotNil(t, beans.NonZero(amount).Validate())
	})

	t.Run("empty amount is not zero", func(t *testing.T) {
		var amount beans.Amount
		assert.Nil(t, beans.NonZero(amount).Validate())
	})
}
