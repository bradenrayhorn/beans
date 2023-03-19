package beans_test

import (
	"math"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	var tests = []struct {
		name     string
		err      string
		expected string
		amounts  []beans.Amount
	}{
		{"add positives", "", "7.6", []beans.Amount{beans.NewAmount(45, -1), beans.NewAmount(31, -1)}},
		{"add mix", "", "1.4", []beans.Amount{beans.NewAmount(45, -1), beans.NewAmount(-31, -1)}},
		{"negative result", "", "-0.5", []beans.Amount{beans.NewAmount(45, -1), beans.NewAmount(-5, 0)}},
		{"single input", "", "1.4", []beans.Amount{beans.NewAmount(14, -1)}},
		{"handles empty amount", "", "0", []beans.Amount{beans.NewEmptyAmount()}},
		{"handles error", "exponent out of range", "", []beans.Amount{beans.NewAmount(math.MaxInt64, math.MaxInt32)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			amount, err := beans.Arithmetic.Add(test.amounts...)

			if len(test.err) > 0 {
				assert.ErrorContains(t, err, test.err)
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.expected, amount.String())
			}
		})
	}
}

func TestNegate(t *testing.T) {
	var tests = []struct {
		name     string
		amount   beans.Amount
		expected string
	}{
		{"can negate negative", beans.NewAmount(-55, -1), "5.5"},
		{"can negate positive", beans.NewAmount(55, -1), "-5.5"},
		{"can negate zero", beans.NewAmount(0, 0), "0"},
		{"can negate empty", beans.NewEmptyAmount(), "0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beforeNegation := test.amount.String()
			negated := beans.Arithmetic.Negate(test.amount)

			assert.Equal(t, test.expected, negated.String())

			assert.Equal(t, beforeNegation, test.amount.String())
		})
	}
}
