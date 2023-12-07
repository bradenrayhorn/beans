package beans

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type ValidatableString string

func (s ValidatableString) Empty() bool {
	return strings.TrimSpace(string(s)) == ""
}

func (s ValidatableString) Length() int {
	return len(string(s))
}

type NullString struct {
	string string
	set    bool
}

func NewNullString(string string) NullString {
	return NullString{string: strings.TrimSpace(string), set: len(strings.TrimSpace(string)) > 0}
}

func (s NullString) SQLNullString() sql.NullString {
	return sql.NullString{String: s.string, Valid: s.set}
}

func NullStringFromSQL(s sql.NullString) NullString {
	return NullString{string: s.String, set: s.Valid}
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.set {
		return json.Marshal(s.string)
	} else {
		return json.Marshal(nil)
	}
}

func (s *NullString) UnmarshalJSON(b []byte) error {
	var string string
	if err := json.Unmarshal(b, &string); err != nil {
		return err
	}
	if len(strings.TrimSpace(string)) > 0 {
		s.string = strings.TrimSpace(string)
		s.set = true
	}
	return nil
}

func (s NullString) Empty() bool {
	if s.set {
		return s.string == ""
	} else {
		return true
	}
}

func (s NullString) Length() int {
	if s.set {
		return len(s.string)
	} else {
		return 0
	}
}

func (s NullString) String() string {
	if s.set {
		return s.string
	} else {
		return ""
	}
}
