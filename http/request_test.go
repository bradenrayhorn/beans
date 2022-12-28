package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeRequest(t *testing.T) {
	request := func(body string) *http.Request {
		return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	}

	t.Run("does not fail with empty body", func(t *testing.T) {
		req := request(``)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		require.Nil(t, err)
	})

	t.Run("returns beans error for bad syntax", func(t *testing.T) {
		req := request(`{"name":}`)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		require.NotNil(t, err)
		var beansError beans.Error
		require.ErrorAs(t, err, &beansError)
		code, _ := beansError.BeansError()
		assert.Equal(t, beans.EUNPROCESSABLE, code)
	})

	t.Run("returns beans error with missing quotes syntax", func(t *testing.T) {
		req := request(`{"name:}`)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		require.NotNil(t, err)
		var beansError beans.Error
		require.ErrorAs(t, err, &beansError)
		code, _ := beansError.BeansError()
		assert.Equal(t, beans.EUNPROCESSABLE, code)
	})

	t.Run("returns beans error for bad type", func(t *testing.T) {
		req := request(`{"name":"test"}`)
		var r struct {
			Name bool `json:"name"`
		}
		err := decodeRequest(req, &r)

		require.NotNil(t, err)
		var beansError beans.Error
		require.ErrorAs(t, err, &beansError)
		code, msg := beansError.BeansError()
		assert.Equal(t, beans.EUNPROCESSABLE, code)
		assert.Equal(t, "Invalid data `string` for `name` field.", msg)
	})
}
