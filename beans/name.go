package beans

import (
	"errors"
	"strings"
)

type Name string

func (n Name) Validate() error {
	trimmedName := strings.TrimSpace(string(n))
	if trimmedName == "" {
		return errors.New("Name is required")
	}
	if len(trimmedName) > 255 {
		return errors.New("Name must be at most 255 characters")
	}
	return nil
}
