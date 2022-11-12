package beans

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeMonth(t *testing.T) {
	date := time.Date(2022, 05, 26, 6, 6, 7, 8, time.Now().Location())

	normalized := NormalizeMonth(date)
	year, month, day := normalized.Date()
	assert.Equal(t, 2022, year)
	assert.Equal(t, time.Month(5), month)
	assert.Equal(t, 1, day)

	hour, minute, second := normalized.Clock()
	assert.Equal(t, 0, hour)
	assert.Equal(t, 0, minute)
	assert.Equal(t, 0, second)
}
