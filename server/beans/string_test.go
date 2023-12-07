package beans

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullStringMarshal(t *testing.T) {
	var tests = []struct {
		name     string
		string   string
		set      bool
		expected string
	}{
		{"not set to null", "", false, `null`},
		{"not set and string to null", "gibberish", false, `null`},
		{"set and empty string", "", true, `""`},
		{"set and full string", "data", true, `"data"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NullString{string: test.string, set: test.set}
			r, err := json.Marshal(s)
			require.Nil(t, err)
			assert.JSONEq(t, test.expected, string(r))
		})
	}
}

func TestNullStringUnmarshal(t *testing.T) {
	var tests = []struct {
		name   string
		json   string
		string string
		set    bool
	}{
		{"null is null", `null`, "", false},
		{"filled string is set", `"data"`, "data", true},
		{"blank string is null", `""`, "", false},
		{"whitespace string is null", `"  "`, "", false},
		{"whitespace is trimmed", `" data  "`, "data", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s NullString
			err := json.Unmarshal([]byte(test.json), &s)
			require.Nil(t, err)
			assert.Equal(t, test.set, s.set)
			assert.Equal(t, test.string, s.string)
		})
	}
}
