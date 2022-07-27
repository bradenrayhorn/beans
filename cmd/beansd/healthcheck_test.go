package main_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r := ta.GetRequest(t, "health-check", nil)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}
