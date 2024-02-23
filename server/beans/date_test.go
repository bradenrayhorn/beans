package beans

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateUnmarshal(t *testing.T) {

	t.Run("can unmarshal valid date", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`"2022-05-11"`), &d)
		require.Nil(t, err)
		assert.Equal(t, "2022-05-11", d.String())
	})

	t.Run("can unmarshal null into empty date", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`null`), &d)
		require.Nil(t, err)
		assert.Equal(t, true, d.Empty())
	})

	t.Run("unmarshal invalid date returns unmarshal error", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`"2022-0511"`), &d)
		require.NotNil(t, err)
		var jsonErr *json.UnmarshalTypeError
		require.ErrorAs(t, err, &jsonErr)
		assert.Equal(t, `"2022-0511"`, jsonErr.Value)
	})
}
