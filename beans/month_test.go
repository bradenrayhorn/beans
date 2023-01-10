package beans

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeMonth(t *testing.T) {
	date := time.Date(2022, 05, 26, 6, 6, 7, 8, time.Now().Location())

	normalized := normalizeMonth(date)
	year, month, day := normalized.Date()
	assert.Equal(t, 2022, year)
	assert.Equal(t, time.Month(5), month)
	assert.Equal(t, 1, day)

	hour, minute, second := normalized.Clock()
	assert.Equal(t, 0, hour)
	assert.Equal(t, 0, minute)
	assert.Equal(t, 0, second)
}

func TestTimeZoneIsStripped(t *testing.T) {
	date1 := NewMonthDate(NewDate(time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)))
	loc, err := time.LoadLocation("America/New_York")
	require.Nil(t, err)

	date2 := NewMonthDate(NewDate(time.Date(2022, 05, 26, 23, 50, 0, 0, loc)))
	date3 := NewMonthDate(NewDate(time.Date(2022, 05, 25, 23, 50, 0, 0, loc)))

	assert.Equal(t, date1, date2)
	assert.Equal(t, date1, date3)
	assert.Equal(t, date2, date3)
}
