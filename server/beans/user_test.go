package beans_test

import (
	"strings"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsernameValidation(t *testing.T) {
	var tests = []struct {
		name     string
		username string
		error    string
	}{
		{"is required", "", "Username is required."},
		{"is required ignores whitespace", "    ", "Username is required."},
		{"can be 32 characters", strings.Repeat("a", 32), ""},
		{"cannot be 33 characters", strings.Repeat("a", 33), "Username must be at most 32 characters."},
		{"can be valid", "user", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			username := beans.Username(test.username)
			err := beans.ValidateFields(username.ValidatableField())
			if test.error == "" {
				assert.Nil(t, err)
			} else {
				require.NotNil(t, err)
				_, msg := err.(beans.Error).BeansError()
				assert.Equal(t, test.error, msg)
			}
		})
	}
}

func TestPasswordValidation(t *testing.T) {
	var tests = []struct {
		name     string
		password string
		error    string
	}{
		{"is required", "", "Password is required."},
		{"is required does not ignores whitespace", "    ", ""},
		{"can be 255 characters", strings.Repeat("a", 255), ""},
		{"cannot be 256 characters", strings.Repeat("a", 256), "Password must be at most 255 characters."},
		{"can be valid", "password", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			password := beans.Password(test.password)
			err := beans.ValidateFields(password.ValidatableField())
			if test.error == "" {
				assert.Nil(t, err)
			} else {
				require.NotNil(t, err)
				_, msg := err.(beans.Error).BeansError()
				assert.Equal(t, test.error, msg)
			}
		})
	}
}
