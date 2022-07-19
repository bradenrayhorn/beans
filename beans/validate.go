package beans

import (
	"errors"
	"strings"
)

type Validatable interface {
	Validate() error
}

func Validate(objects ...Validatable) error {
	messages := []string{}
	for _, o := range objects {
		if err := o.Validate(); err != nil {
			messages = append(messages, err.Error()+".")
		}
	}

	if len(messages) == 0 {
		return nil
	}

	return WrapError(errors.New(strings.Join(messages, " ")), ErrorInvalid)
}
