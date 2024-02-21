package beans

import (
	"encoding/json"
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

func TestMonthDate(t *testing.T) {

	t.Run("get previous", func(t *testing.T) {
		monthDate := NewMonthDate(NewDate(time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)))

		previous := monthDate.Previous()
		assert.Equal(t, previous.String(), "2022-04-01")
	})
}

func TestMonthDateJSON(t *testing.T) {

	t.Run("can unmarshal date", func(t *testing.T) {
		var d MonthDate
		err := json.Unmarshal([]byte(`"2022-05-11"`), &d)
		require.Nil(t, err)
		assert.Equal(t, "2022-05-01", d.String())
	})

	t.Run("unmarshal invalid date returns unmarshal error", func(t *testing.T) {
		var d MonthDate
		err := json.Unmarshal([]byte(`"2022-0511"`), &d)
		require.NotNil(t, err)
		var jsonErr *json.UnmarshalTypeError
		require.ErrorAs(t, err, &jsonErr)
		assert.Equal(t, `"2022-0511"`, jsonErr.Value)
	})

	t.Run("can marshal date", func(t *testing.T) {
		monthDate := NewMonthDate(NewDate(time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)))

		res, err := json.Marshal(monthDate)
		require.NoError(t, err)
		assert.Equal(t, `"2022-05-01"`, string(res))
	})

	t.Run("can marshal empty date", func(t *testing.T) {
		res, err := json.Marshal(MonthDate{})
		require.NoError(t, err)
		assert.Equal(t, `null`, string(res))
	})
}
